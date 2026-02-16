package replays

import (
	"perdoon/src/replay"
	"perdoon/src/state"
)

type ReplayMaker func(state *state.State) replay.Replay

var REGISTRY map[string]ReplayMaker

func init() {
	REGISTRY = make(map[string]ReplayMaker)

	REGISTRY["echo"] = ReplayMaker(NewEcho)
	REGISTRY["random"] = ReplayMaker(NewRandom)
	REGISTRY["zero"] = ReplayMaker(NewZero)
	REGISTRY["bytes"] = ReplayMaker(NewBytes)
	REGISTRY["potato"] = ReplayMaker(NewPotato)
	REGISTRY["qt"] = ReplayMaker(NewQt)
}
