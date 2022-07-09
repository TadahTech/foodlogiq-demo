package service

import (
	"encoding/json"
	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"net/http"
	"time"

	v "github.com/go-ozzo/ozzo-validation"
)

func (rs *RestServer) createEvent(w http.ResponseWriter, r *http.Request) {
	var event *model.Event
	err := decodeRequest(r, &event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if valErr := v.ValidateStruct(&event, v.Field(&event.Type, v.Required, v.In("shipping", "receiving"))); valErr != nil {
		http.Error(w, "invalid type", http.StatusBadRequest)
		return
	}

	if valErr := v.Validate(&event.Contents, v.By(validateContents)); valErr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// We can ignore the error here because middleware will catch any errors
	user, _ := userFromBearer(r)

	event.CreatedBy = user.UserID
	event.CreatedAt = time.Now().Format(time.RFC3339)

	err = rs.db.CreateEvent(event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResp, _ := json.Marshal(map[string]interface{}{"event_id": event.ID})
	w.Write(jsonResp)
}
