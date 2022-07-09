package service

import (
	"encoding/json"
	"net/http"
)

func decodeRequest(r *http.Request, result interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	err := decoder.Decode(&result)

	return err
}
