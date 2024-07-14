package handlers

import "log"

func HandleErrors(err error) {
	if err != nil {
		log.Panic(err)
	}
}
