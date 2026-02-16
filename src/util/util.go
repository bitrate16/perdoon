package util

import "log"

func WrapErrorLog(prefix string, err error) {
	if err != nil {
		log.Printf("%s: %s", err)
	}
}
