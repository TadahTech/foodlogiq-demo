package data

import (
	"context"
	"fmt"
	"time"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
)

func (d dataStore) GetAllEvents(owner int) ([]*model.Event, error) {
	var results []*model.Event
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "[Mongo] GetAllEvents")
	defer span.Finish()

	resp, err := d.collection.Find(ctx, bson.M{"created_by": owner, "is_deleted": false})

	if err != nil {
		return nil, fmt.Errorf("[Mongo] GetAllEvents for owner %v error: %w", owner, err)
	}

	for resp.Next(context.TODO()) {
		var elem *storedEvent
		err := resp.Decode(&elem)

		if err != nil {
			return nil, fmt.Errorf("[Mongo] GetAllEvents#decodeEvent error: %w", err)
		}

		results = append(results, createEventFromStored(elem))
	}

	return results, nil
}
