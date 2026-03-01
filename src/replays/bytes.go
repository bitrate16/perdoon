package replays

import (
	"encoding/hex"
	"errors"
	"net"
	"perdoon/src/config"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Bytes struct {
	state  *state.State
	bytes  []byte
	config *config.ResponseConfig
}

func NewBytes(
	state *state.State,
	config *config.ResponseConfig,
) replay.Replay {
	return &Bytes{
		state:  state,
		config: config,
	}
}

func (e *Bytes) Init() error {
	var err error

	if len(e.config.Bytes)&0b1 != 0 {
		e.bytes, err = hex.DecodeString("0" + e.config.Bytes)
	} else {
		e.bytes, err = hex.DecodeString(e.config.Bytes)
	}

	if err != nil {
		return err
	}

	if len(e.bytes) == 0 {
		return errors.New("cannot use empty bytes template for bytes strategy")
	}

	return nil
}

func (e *Bytes) Replay(
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
