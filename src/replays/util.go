package replays

import (
	"math/rand"
	"perdoon/src/config"
)

func SelectRandomFromRanges(objects []*config.ValueRange) int {
	if len(objects) == 0 {
		return 0
	}

	randomIndex := rand.Intn(len(objects))
	selectedObject := objects[randomIndex]

	randomNumber := rand.Intn(selectedObject.End-selectedObject.Start+1) + selectedObject.Start

	return randomNumber
}
