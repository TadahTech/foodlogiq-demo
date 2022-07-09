package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func (rs *RestServer) listEvents(w http.ResponseWriter, r *http.Request) {
	user, _ := userFromBearer(r)
	events, err := rs.db.GetAllEvents(user.UserID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "cannot find an event id and created_by match", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("cannot get event %v", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonVal, _ := json.Marshal(events)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonVal)
}
