package state

import (
	"perdoon/src/config"
	"perdoon/src/replay"
	"perdoon/src/track"

	"github.com/gofrs/uuid/v5"
)

type State struct {
	Config  *config.Config
	Track   track.Track
	Replay  replay.Replay
	Session uuid.UUID
}
