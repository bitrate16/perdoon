package tracks

import (
	"perdoon/src/state"
	"perdoon/src/track"
)

type TrackMaker func(state *state.State) track.Track

var REGISTRY map[string]TrackMaker

func init() {
	REGISTRY = make(map[string]TrackMaker)

	REGISTRY["print"] = TrackMaker(NewPrint)
	REGISTRY["sqlite"] = TrackMaker(NewSQLite)
}
