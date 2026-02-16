package server

import "perdoon/src/util"

func WrapTrackErrorLog(err error) {
	util.WrapErrorLog("Failed track event", err)
}
