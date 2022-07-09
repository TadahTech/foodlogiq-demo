package data

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d dataStore) DeleteEvent(eventId string, createdBy int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "[Mongo] DeleteEvent")
	defer span.Finish()

	update := bson.M{
		"$set": bson.D{{"is_deleted", true}},
	}

	filter := bson.M{"_id": eventId, "created_by": createdBy}
	resp, err := d.collection.UpdateOne(ctx, filter, update, options.Update())

	if err != nil {
		return fmt.Errorf("[Mongo] DeleteEvent error: %w", err)
	}

	// If there were 0 matches, return a no document error
	// otherwise, return no error (success)
	if resp.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
