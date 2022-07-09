package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	log "github.com/sirupsen/logrus"
)

var (
	client  = http.Client{}
	eventId = "62c9f05a8fc7c13583b1dc06"
	token   = "Bearer 74edf612f393b4eb01fbc2c29dd96671"
)

func main() {
	initLogging("debug")
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-= GET EVENT -=-=-=-=-=-=-=-=-=-=-=")
	getEvent()
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-= LIST EVENTS -=-=-=-=-=-=-=-=-=-=-=")
	getEvents()
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-= DELETE EVENT -=-=-=-=-=-=-=-=-=-=-=")
	deleteEvent()
}

func getEvent() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/event?event_id="+eventId, nil)

	if err != nil {
		log.WithError(err).Error("failed to make request")
		return
	}

	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("failed to exec request")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed to get event")
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.WithError(err).Error("failed to read all from resp body")
			return
		}

		fmt.Printf("Response body: %v", string(body))
		return
	}

	var mapper *model.Event
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&mapper)

	if err != nil {
		log.WithError(err).Error("failed to decode request")
		return
	}

	jsonVal, _ := json.Marshal(mapper)
	fmt.Println(string(jsonVal))
}

func getEvents() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/event/all", nil)

	if err != nil {
		log.WithError(err).Error("failed to make request")
		return
	}

	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("failed to exec request")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed to get all events")
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.WithError(err).Error("failed to read all from resp body")
			return
		}

		fmt.Printf("Response body: %v", string(body))
		return
	}

	var mapper []*model.Event
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&mapper)

	if err != nil {
		log.WithError(err).Error("failed to decode request")
		return
	}

	jsonVal, _ := json.Marshal(mapper)
	fmt.Println(len(mapper))
	fmt.Println(string(jsonVal))
}

func deleteEvent() {
	request := &model.Event{
		ID: eventId,
	}

	jsonBlock, _ := json.Marshal(&request)
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8000/event", bytes.NewBuffer(jsonBlock))

	if err != nil {
		log.WithError(err).Error("failed to make request")
		return
	}

	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("failed to exec request")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed to delete event")
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.WithError(err).Error("failed to read all from resp body")
			return
		}

		fmt.Printf("Response body: %v", string(body))
		return
	}

	fmt.Println("Deleted event")
}
func createEvent() {
	request := &model.Event{
		Type: "shipping",
		Contents: []*model.Content{
			{
				Gtin: "1234",
				Lot:  "abcdef",
			},
		},
	}

	jsonBlock, _ := json.Marshal(&request)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/event", bytes.NewBuffer(jsonBlock))

	if err != nil {
		log.WithError(err).Error("failed to make request")
		return
	}

	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Error("failed to exec request")
		return
	}

	if resp.StatusCode != http.StatusCreated {
		log.Info("failed to create event")
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.WithError(err).Error("failed to read all from resp body")
			return
		}

		log.Infof("Response body: %v", string(body))
		return
	}

	var mapper map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&mapper)

	if err != nil {
		log.WithError(err).Error("could not decode response")
		return
	}

	eventId = mapper["event_id"].(string)

	log.Infof("Made event %v", eventId)
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
