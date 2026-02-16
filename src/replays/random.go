package replays

import (
	"log"
	"math/rand"
	"net"
	"perdoon/src/replay"
	"perdoon/src/state"
)

type Random struct {
	state *state.State
}

func NewRandom(state *state.State) replay.Replay {
	return &Random{
		state: state,
	}
}

func (e *Random) Init() error {
	return nil
}

func (e *Random) Replay(
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
	_, err := rand.Read(buffer)
	if err != nil {
		log.Printf("Failed generate random bytes: %s", err)
		return []byte{}
	}

	return buffer
}
