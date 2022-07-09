package main

import (
	"bytes"
	"encoding/json"
	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func main() {
	client := http.Client{}
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
		logrus.WithError(err).Error("failed to make request")
		return
	}

	req.Header.Add("Authorization", "Bearer 74edf612f393b4eb01fbc2c29dd96671")

	resp, err := client.Do(req)

	if err != nil {
		logrus.WithError(err).Error("failed to exec request")
		return
	}

	if resp.StatusCode != http.StatusCreated {
		logrus.Info("failed to create event")
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			logrus.WithError(err).Error("failed to read all from resp body")
			return
		}

		logrus.Infof("Response body: %v", string(body))
		return
	}

	var mapper map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&mapper)

	if err != nil {
		logrus.WithError(err).Error("could not decode response")
		return
	}
	logrus.Infof("Made event %v", mapper["event_id"])
}
