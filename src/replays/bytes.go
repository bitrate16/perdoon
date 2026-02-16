package replays

import (
	"encoding/hex"
	"net"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Bytes struct {
	state *state.State
	bytes []byte
}

func NewBytes(state *state.State) replay.Replay {
	return &Bytes{
		state: state,
	}
}

func (e *Bytes) Init() error {
	var err error

	if len(e.state.Config.Response.Bytes)&0b1 != 0 {
		e.bytes, err = hex.DecodeString("0" + e.state.Config.Response.Bytes)
	} else {
		e.bytes, err = hex.DecodeString(e.state.Config.Response.Bytes)
	}

	return err
}

func (e *Bytes) Replay(
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

		if diff > len(e.bytes) {
			copy(buffer[offset:offset+len(e.bytes)], e.bytes)
			offset += len(e.bytes)
		} else {
			copy(buffer[offset:offset+diff], e.bytes[:diff])
			break
		}
	}

	return buffer
}
