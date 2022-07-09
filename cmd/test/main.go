package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	client  = http.Client{}
	eventId string
	token   = "Bearer 74edf612f393b4eb01fbc2c29dd96671"
)

func main() {
	initLogging("debug")
	getEvent()
	getEvents()
	deleteEvent()
}

func initLogging(level string) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "min":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
