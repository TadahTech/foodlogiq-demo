package main

import (
	"net/http"
	"os"

	"github.com/TadahTech/foodlogiq-demo/pkg/data"
	"github.com/TadahTech/foodlogiq-demo/pkg/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	initLogging("info")
	connectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if len(connectionString) == 0 {
		log.Error("connection string is empty")
		return
	}

	log.Info("Setting up data store")
	db, err := data.NewDataStore(connectionString)

	if err != nil {
		log.WithError(err).Error("failed to create database")
		return
	}

	log.Info("Setting up rest server")
	s := service.NewServer(db)
	log.Info("Serving HTTP server on port 8000")
	if err := http.ListenAndServe(":8000", s.Router); err != nil {
		log.WithError(err).Error("server error")
		return
	}
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
