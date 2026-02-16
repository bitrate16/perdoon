package replays

import (
	"net"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Uwu struct {
	state *state.State
}

const QT = ":3"

func NewQt(state *state.State) replay.Replay {
	return &Uwu{
		state: state,
	}
}

func (e *Uwu) Init() error {
	return nil
}

func (e *Uwu) Replay(
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
