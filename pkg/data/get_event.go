package data

import (
	"context"
	"fmt"
	"time"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d dataStore) GetEvent(eventId string, createdBy int) (*model.Event, error) {
	var event *storedEvent
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "[Mongo] GetEvent")
	defer span.Finish()

	filter := bson.M{"_id": eventId, "created_by": createdBy, "is_deleted": false}

	err := d.collection.FindOne(ctx, filter).Decode(&event)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no documents were found matching that eventId and created_by")
		}
		return nil, fmt.Errorf("[Mongo] GetEvent error: %w", err)
	}

	return createEventFromStored(event), nil
}
