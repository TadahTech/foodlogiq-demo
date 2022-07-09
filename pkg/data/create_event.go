package data

import (
	"context"
	"fmt"
	"time"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (d dataStore) CreateEvent(event *model.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	span, ctx := opentracing.StartSpanFromContext(ctx, "[Mongo] CreateEvent")
	defer span.Finish()

	sc, err := createStoredEvent(event)

	if err != nil {
		return fmt.Errorf("[Mongo] CreateEvent#createStoredEvent error: %w", err)
	}

	_, err = d.collection.InsertOne(ctx, sc)

	if err != nil {
		return fmt.Errorf("[Mongo] CreateEvent#insert error: %w", err)
	}

	event.ID = sc.ID.String()

	return nil
}
