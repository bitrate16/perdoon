// perdoon - port sniffer
// Copyright (C) 2026  bitrate16
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"perdoon/src/config"
	"perdoon/src/replays"
	"perdoon/src/server"
	"perdoon/src/state"
	"perdoon/src/tracks"
	"reflect"

	"github.com/gofrs/uuid/v5"
)

func main() {
	// State prepare
	state := &state.State{}
	var err error

	// Session ID
	session, err := uuid.NewV4()
	state.Session = session

	// Config parse
	state.Config, err = config.LoadConfig(config.GetConfigPath())
	if err != nil {
		log.Printf("Startup failed: %s", err)
		return
	}

	// Replayer init for TCP
	if tcpReplayMaker, ok := replays.REGISTRY[state.Config.TCP.Response.Strategy]; ok {
		state.TCPReplay = tcpReplayMaker(state, state.Config.TCP.Response)
		err := state.TCPReplay.Init()
		if err != nil {
			log.Printf("TCP strategy Init failed: %s", err)
			return
		}
	} else {
		log.Printf("TCP strategy: %s not supported", state.Config.TCP.Response.Strategy)
		os.Exit(1)
	}

	// Replayer init for UDP
	if udpReplayMaker, ok := replays.REGISTRY[state.Config.UDP.Response.Strategy]; ok {
		state.UDPReplay = udpReplayMaker(state, state.Config.UDP.Response)
		err := state.UDPReplay.Init()
		if err != nil {
			log.Printf("UDP strategy Init failed: %s", err)
			return
		}
	} else {
		log.Printf("UDP strategy: %s not supported", state.Config.UDP.Response.Strategy)
		os.Exit(1)
	}

	if state.Config.Debug {
		state.Track = tracks.NewPrint(state)
	} else {
		state.Track = tracks.NewSQLite(state)
	}

	err = state.Track.Open()
	if err != nil {
		log.Printf("Failed to prepare track: %s", err)
	}
	defer state.Track.Close()

	// Start servers
	udpServers := make(map[int]server.ServerController)
	tcpServers := make(map[int]server.ServerController)

	// Exclude ports
	udpExclude := make(map[int]struct{})
	tcpExclude := make(map[int]struct{})

	// Build exclude maps
	for _, port := range state.Config.UDP.Exclude {
		udpExclude[port] = struct{}{}
	}
	for _, port := range state.Config.TCP.Exclude {
		tcpExclude[port] = struct{}{}
	}

	// UDP Servers
	for _, portRange := range state.Config.UDP.Ports {
		for port := portRange.Start; port <= portRange.End; port++ {
			if port <= 0 || port >= 65535 {
				continue
			}

			// Skip excluded
			if _, skip := udpExclude[port]; skip {
				continue
			}

			// Skip existing
			if _, has := udpServers[port]; has {
				continue
			}

			// Create server
			udpServer := server.NewUDPServer(state, port)

			// Start server
			err := udpServer.Start()
			if err != nil {
				udpServer = nil
				log.Printf("Failed start UDP server on port %d: %s", port, err)
			}

			// Add anyway
			udpServers[port] = udpServer
		}
	}

	// TCP Servers
	for _, portRange := range state.Config.TCP.Ports {
		for port := portRange.Start; port <= portRange.End; port++ {
			if port <= 0 || port >= 65535 {
				continue
			}

			// Skip excluded
			if _, skip := tcpExclude[port]; skip {
				continue
			}

			// Skip existing
			if _, has := tcpServers[port]; has {
				continue
			}

			// Create server
			tcpServer := server.NewTCPServer(state, port)

			// Start server
			err := tcpServer.Start()
			if err != nil {
				tcpServer = nil
				log.Printf("Failed start TCP server on port %d: %s", port, err)
			}

			// Add anyway
			tcpServers[port] = tcpServer
		}
	}

	fmt.Println("Service running")

	// Run forever until quit
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan

	fmt.Println("Interrupted")

	// Terminate UDP
	for _, serverController := range udpServers {
		if serverController == nil || reflect.ValueOf(serverController).Kind() == reflect.Ptr && reflect.ValueOf(serverController).IsNil() {
			continue
		}

		serverController.Stop()
	}

	// Terminate TCP
	for _, serverController := range tcpServers {
		if serverController == nil || reflect.ValueOf(serverController).Kind() == reflect.Ptr && reflect.ValueOf(serverController).IsNil() {
			continue
		}

		serverController.Stop()
	}
}
