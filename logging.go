package main

import (
	log "github.com/sirupsen/logrus"
)

func setupLogging(enableDebug bool) error {
	if enableDebug {
		engine.EnableQueryDebug()
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return nil
}