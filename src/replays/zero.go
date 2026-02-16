package replays

import (
	"net"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Zero struct {
	state *state.State
}

func NewZero(state *state.State) replay.Replay {
	return &Zero{
		state: state,
	}
}

func (e *Zero) Init() error {
	return nil
}

func (e *Zero) Replay(
	client_payload []byte,
	protocol string,
	port int,
	clientIP net.IP,
	clientPort int,
) []byte {
	size := SelectRandomFromRanges(e.state.Config.Response.Sizes)
	if size <= 0 {
		return []byte{}
	}

	return make([]byte, size)
}
