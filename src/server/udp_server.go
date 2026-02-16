package server

import (
	"fmt"
	"log"
	"net"
	"perdoon/src/state"

	"github.com/gofrs/uuid/v5"
)

type UDPServer struct {
	state *state.State
	port  int

	conn *net.UDPConn

	open     bool
	exitChan chan struct{}
}

func NewUDPServer(
	state *state.State,
	port int,
) *UDPServer {
	return &UDPServer{
		state:    state,
		port:     port,
		exitChan: make(chan struct{}),
	}
}

func (s *UDPServer) Start() error {
	s.exitChan = make(chan struct{})
	s.open = true

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("Error resolving UDP address: %s", err)
	}

	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("Error listening on UDP: %s", err)
	}

	// Start reader
	go func() {
		// Prepare buffer for incoming data
		buffer := make([]byte, s.state.Config.UDP.ChunkSize)

		// Receive data from clients
		for s.open {
			chainId, err := uuid.NewV4()

			// Read incoming message
			dataSize, clientAddr, err := s.conn.ReadFromUDP(buffer)

			if !s.open {
				break
			}

			// Track message
			WrapTrackErrorLog(s.state.Track.Event(
				"udp",
				s.port,
				clientAddr.IP,
				clientAddr.Port,
				chainId,
				"connect",
				0,
				nil,
			))

			// Track error
			if err != nil {
				log.Printf("UDP READ ERROR [PORT: %d] [CLIENT: %s]: %s", s.port, clientAddr, err)

				// Track error
				WrapTrackErrorLog(s.state.Track.Event(
					"udp",
					s.port,
					clientAddr.IP,
					clientAddr.Port,
					chainId,
					"read_error",
					0,
					nil,
				))

				continue
			}

			// Track client payload
			if s.state.Config.Database.Record.RequestPayload {
				WrapTrackErrorLog(s.state.Track.Event(
					"udp",
					s.port,
					clientAddr.IP,
					clientAddr.Port,
					chainId,
					"message",
					dataSize,
					buffer[:dataSize],
				))
			} else {
				WrapTrackErrorLog(s.state.Track.Event(
					"udp",
					s.port,
					clientAddr.IP,
					clientAddr.Port,
					chainId,
					"message",
					dataSize,
					nil,
				))
			}

			response := s.state.Replay.Replay(
				buffer[:dataSize],
				"udp",
				s.port,
				clientAddr.IP,
				clientAddr.Port,
			)

			// Track response
			if s.state.Config.Database.Record.ResponsePayload {
				WrapTrackErrorLog(s.state.Track.Event(
					"udp",
					s.port,
					clientAddr.IP,
					clientAddr.Port,
					chainId,
					"reponse",
					dataSize,
					response,
				))
			} else {
				WrapTrackErrorLog(s.state.Track.Event(
					"udp",
					s.port,
					clientAddr.IP,
					clientAddr.Port,
					chainId,
					"reponse",
					dataSize,
					nil,
				))
			}

			_, err = s.conn.WriteToUDP(response, clientAddr)
			if err != nil {
				log.Printf("UDP WRITE ERROR [PORT: %d] [CLIENT: %s]: %s", s.port, clientAddr, err)
			}
		}

		// Mark stopped
		s.exitChan <- struct{}{}
	}()

	return nil
}

func (s *UDPServer) Stop() {
	s.open = false
	s.conn.Close()

	// Wait for stop
	<-s.exitChan
}
