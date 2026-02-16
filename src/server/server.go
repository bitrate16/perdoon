package server

import (
	"fmt"
	"perdoon/src/state"
)

// Interface used to gracefully close listening servers
type ServerController interface {
	Start() error
	Stop()
}

func NewServer(
	state *state.State,
	protocol string,
	port int,
) (ServerController, error) {
	if protocol == "tcp" {
		return NewTCPServer(state, port), nil
	}

	if protocol == "udp" {
		return NewUDPServer(state, port), nil
	}

	return nil, fmt.Errorf("protocol %s not supported", protocol)
}
