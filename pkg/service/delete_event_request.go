package service

import (
	"fmt"
	"net/http"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func (rs *RestServer) deleteEvent(w http.ResponseWriter, r *http.Request) {
	var event *model.Event
	err := decodeRequest(r, &event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(event.ID) == 0 {
		http.Error(w, "event id is empty", http.StatusBadRequest)
		return
	}

	user, _ := userFromBearer(r)
	err = rs.db.DeleteEvent(event.ID, user.UserID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "cannot find an event id and created_by match", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("cannot delete event %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
