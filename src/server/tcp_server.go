package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"perdoon/src/state"

	"github.com/gofrs/uuid/v5"
)

type TCPServer struct {
	state *state.State
	port  int

	conn *net.TCPListener

	open     bool
	exitChan chan struct{}
}

func NewTCPServer(
	state *state.State,
	port int,
) *TCPServer {
	return &TCPServer{
		state:    state,
		port:     port,
		exitChan: make(chan struct{}),
	}
}

func (s *TCPServer) Start() error {
	s.exitChan = make(chan struct{})
	s.open = true

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("Error resolving TCP address: %s", err)
	}

	s.conn, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("Error listening on TCP: %s", err)
	}

	// Start reader
	go func() {
		// Receive data from clients
		for s.open {
			chainId, err := uuid.NewV4()

			// Read incoming message
			// dataSize, clientAddr, err := s.conn.AcceptTCP()
			conn, err := s.conn.AcceptTCP()
			if !s.open {
				break
			}
			clientAddr := conn.RemoteAddr().(*net.TCPAddr)

			if err != nil {
				// Track error
				WrapTrackErrorLog(s.state.Track.Event(
					"tcp",
					s.port,
					clientAddr.IP,
					clientAddr.Port,
					chainId,
					"connect_error",
					0,
					nil,
				))

				break
			}
			defer conn.Close()

			// Track message
			WrapTrackErrorLog(s.state.Track.Event(
				"tcp",
				s.port,
				clientAddr.IP,
				clientAddr.Port,
				chainId,
				"connect",
				0,
				nil,
			))

			if !s.open {
				break
			}

			go func() {
				// Prepare buffer for incoming data
				buffer := make([]byte, s.state.Config.TCP.ChunkSize)

				for s.open {
					// Read incoming message
					dataSize, err := conn.Read(buffer)

					if !s.open {
						break
					}

					// Track error
					if err != nil {

						// Client disconnect
						if err == io.EOF {

							// Track error
							WrapTrackErrorLog(s.state.Track.Event(
								"tcp",
								s.port,
								clientAddr.IP,
								clientAddr.Port,
								chainId,
								"disconnect",
								0,
								nil,
							))

							return
						}

						log.Printf("TCP READ ERROR [PORT: %d] [CLIENT: %s]: %s", s.port, clientAddr, err)

						// Track error
						WrapTrackErrorLog(s.state.Track.Event(
							"tcp",
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
							"tcp",
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
							"tcp",
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
						"tcp",
						s.port,
						clientAddr.IP,
						clientAddr.Port,
					)

					// Track response
					if s.state.Config.Database.Record.ResponsePayload {
						WrapTrackErrorLog(s.state.Track.Event(
							"tcp",
							s.port,
							clientAddr.IP,
							clientAddr.Port,
							chainId,
							"reponse",
							len(response),
							response,
						))
					} else {
						WrapTrackErrorLog(s.state.Track.Event(
							"tcp",
							s.port,
							clientAddr.IP,
							clientAddr.Port,
							chainId,
							"reponse",
							len(response),
							nil,
						))
					}

					_, err = conn.Write(response)
					if err != nil {
						log.Printf("TCP WRITE ERROR [PORT: %d] [CLIENT: %s]: %s", s.port, clientAddr, err)
					}
				}
			}()
		}

		// Mark stopped
		s.exitChan <- struct{}{}
	}()

	return nil
}

func (s *TCPServer) Stop() {
	s.open = false
	s.conn.Close()

	// Wait for stop
	<-s.exitChan
}
