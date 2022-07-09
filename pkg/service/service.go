package service

import (
	"net/http"

	"github.com/TadahTech/foodlogiq-demo/pkg/data"
	"github.com/gorilla/mux"
)

type RestServer struct {
	db     data.EventsMongoDB
	Router *mux.Router
}

func NewServer(db data.EventsMongoDB) *RestServer {
	r := mux.NewRouter()

	server := &RestServer{
		db:     db,
		Router: r,
	}

	r.HandleFunc("/health", healthCheck)
	r.HandleFunc("/event", server.createEvent).Methods(http.MethodPost)
	r.HandleFunc("/event", server.getEvent).Methods(http.MethodGet)
	r.HandleFunc("/event/all", server.listEvents).Methods(http.MethodGet)
	r.HandleFunc("/event", server.deleteEvent).Methods(http.MethodDelete)

	r.Use(tokenMw)

	return server
}
