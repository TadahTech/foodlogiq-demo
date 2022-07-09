package data

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d dataStore) DeleteEvent(eventId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "[Mongo] DeleteEvent")
	defer span.Finish()

	update := bson.M{
		"$set": bson.D{{"is_deleted", "true"}},
	}

	filter := bson.M{"_id": eventId}
	_, err := d.collection.UpdateOne(ctx, filter, update, options.Update())

	if err != nil {
		return fmt.Errorf("[Mongo] DeleteEvent error: %w", err)
	}

	return nil
}
