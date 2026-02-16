package replays

import (
	"net"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Echo struct {
}

func NewEcho(state *state.State) replay.Replay {
	return &Echo{}
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
