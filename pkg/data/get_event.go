package data

import (
	"context"
	"fmt"
	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (d dataStore) GetEvent(eventId string) (*model.Event, error) {
	var event *storedEvent
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "[Mongo] GetEvent")
	defer span.Finish()

	filter := bson.M{"id": eventId}

	err := d.collection.FindOne(ctx, filter).Decode(&event)

	if err != nil {
		return nil, fmt.Errorf("[Mongo] GetEvent error: %w", err)
	}

	return createEventFromStored(event), nil
}
