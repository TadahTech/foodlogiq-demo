package service

import (
	"encoding/json"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(map[string]interface{}{"status": "ok"})
	w.Write(jsonResp)
}
