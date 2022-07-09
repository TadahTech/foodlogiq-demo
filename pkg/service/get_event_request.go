package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func (rs *RestServer) getEvent(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	eventId := v.Get("event_id")

	if len(eventId) == 0 {
		http.Error(w, "event id is empty", http.StatusBadRequest)
		return
	}

	user, _ := userFromBearer(r)
	event, err := rs.db.GetEvent(eventId, user.UserID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "cannot find an event id and created_by match", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("cannot get event %v", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonVal, _ := json.Marshal(event)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonVal)
}
