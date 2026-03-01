package replays

import (
	"net"
	"perdoon/src/config"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Potato struct {
	state  *state.State
	config *config.ResponseConfig
}

const POTATO = "potato"

func NewPotato(
	state *state.State,
	config *config.ResponseConfig,
) replay.Replay {
	return &Potato{
		state:  state,
		config: config,
	}
}

func (e *Potato) Init() error {
	return nil
}

func (e *Potato) Replay(
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

	buffer := make([]byte, size)
	offset := 0

	for {
		diff := size - offset

		if diff > len(POTATO) {
			copy(buffer[offset:offset+len(POTATO)], POTATO)
			offset += len(POTATO)
		} else {
			copy(buffer[offset:offset+diff], POTATO[:diff])
			break
		}
	}

	return buffer
}
