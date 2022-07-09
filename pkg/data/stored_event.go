package data

import (
	"fmt"
	"time"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// storedEvent BSON representation of a JSON Event with primitive types and ID
type storedEvent struct {
	CreatedAt time.Time          `bson:"created_at"`
	IsDeleted bool               `bson:"is_deleted"`
	CreatedBy int                `bson:"created_by"`
	Contents  []*storedContent   `bson:"contents"`
	ID        primitive.ObjectID `bson:"_id"`
	Type      string             `bson:"type"`
}

type storedContent struct {
	Lot            string     `json:"lot"`
	Gtin           string     `bson:"gtin"`
	BestByDate     *time.Time `bson:"best_by_date"`
	ExpirationDate *time.Time `bson:"expiration_date"`
}

// createStoredEvent Convert a JSON Event into a BSON event and keep the fields as primitive to Mongo as possible for betting filtering
func createStoredEvent(e *model.Event) (*storedEvent, error) {
	event := &storedEvent{
		CreatedAt: time.Now(),
		IsDeleted: false,
		CreatedBy: e.CreatedBy,
		ID:        primitive.NewObjectID(),
		Type:      e.Type,
	}

	var storedContents []*storedContent

	for _, c := range e.Contents {
		sc := &storedContent{
			Lot:  c.Lot,
			Gtin: c.Gtin,
		}

		if len(c.BestByDate) > 0 {
			t, err := time.Parse(time.RFC3339, c.BestByDate)
			if err != nil {
				return nil, fmt.Errorf("parsing bestByDate error %v", err)
			}
			sc.BestByDate = &t
		}

		if len(c.ExpirationDate) > 0 {
			t, err := time.Parse(time.RFC3339, c.ExpirationDate)
			if err != nil {
				return nil, fmt.Errorf("parsing expirationDate error %v", err)
			}
			sc.ExpirationDate = &t
		}

		storedContents = append(storedContents, sc)
	}

	event.Contents = storedContents

	return event, nil

}

func createEventFromStored(e *storedEvent) *model.Event {
	event := &model.Event{
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
		IsDeleted: e.IsDeleted,
		CreatedBy: e.CreatedBy,
		ID:        e.ID.Hex(),
		Type:      e.Type,
	}

	var contents []*model.Content

	for _, c := range e.Contents {
		sc := &model.Content{
			Lot:  c.Lot,
			Gtin: c.Gtin,
		}

		if c.BestByDate != nil {
			sc.BestByDate = c.BestByDate.Format(time.RFC3339)
		}

		if c.ExpirationDate != nil {
			sc.ExpirationDate = c.ExpirationDate.Format(time.RFC3339)
		}

		contents = append(contents, sc)
	}

	event.Contents = contents
	return event
}
