package replays

import (
	"net"
	"perdoon/src/config"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Echo struct {
	config *config.ResponseConfig
}

func NewEcho(
	state *state.State,
	config *config.ResponseConfig,
) replay.Replay {
	return &Echo{
		config: config,
	}
}

func (e *Echo) Init() error {
	return nil
}

func (e *Echo) Replay(
	client_payload []byte,
	protocol string,
	port int,
	clientIP net.IP,
	clientPort int,
) []byte {
	return client_payload
}
