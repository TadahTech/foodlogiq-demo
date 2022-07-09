package data

import (
	"context"
	"errors"
	"time"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "foodlogiq"
	collectionName = "events"
)

type EventsMongoDB interface {
	CreateEvent(event *model.Event) (string, error)
	DeleteEvent(eventId string, createdBy int) error
	GetEvent(eventId string, createdBy int) (*model.Event, error)
	GetAllEvents(owner string) ([]*model.Event, error)
}

type dataStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type EventsDatastore struct {
	EventsMongoDB
}

func NewDataStore(connection string) (*EventsDatastore, error) {
	newClient, err := mongo.NewClient(options.Client().ApplyURI(connection))

	if err != nil {
		log.WithError(err).Error("could not create mongo client")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = newClient.Connect(ctx)
	if err != nil {
		log.WithError(err).Error("could not connect to mongo")
		return nil, err
	}

	client := newClient

	collection := client.Database(dbName).Collection(collectionName)

	if collection == nil {
		return nil, errors.New("could not connect to the collection")
	}

	eventsDB := &dataStore{
		client:     client,
		collection: collection,
	}

	return &EventsDatastore{eventsDB}, nil
}
