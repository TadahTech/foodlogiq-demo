package main

import (
	"github.com/TadahTech/foodlogiq-demo/pkg/data"
	"github.com/TadahTech/foodlogiq-demo/pkg/service"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		return
	}

	connectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if len(connectionString) == 0 {
		log.Error("connection string is empty")
		return
	}

	db, err := data.NewDataStore(connectionString)

	if err != nil {
		log.WithError(err).Error("failed to create database")
		return
	}

	s := service.NewServer(db)
	if err := http.ListenAndServe("localhost:8000", s.Router); err != nil {
		log.WithError(err).Error("server error")
		return
	}
}
