package replays

import (
	"net"
	"perdoon/src/config"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Zero struct {
	state  *state.State
	config *config.ResponseConfig
}

func NewZero(
	state *state.State,
	config *config.ResponseConfig,
) replay.Replay {
	return &Zero{
		state:  state,
		config: config,
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
	size := SelectRandomFromRanges(e.config.Sizes)
	if size <= 0 {
		return []byte{}
	}

	return make([]byte, size)
}
