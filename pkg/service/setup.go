package service

import (
	"net/http"

	"github.com/TadahTech/foodlogiq-demo/pkg/data"
	"github.com/gorilla/mux"
)

type RestServer struct {
	db     data.EventsMongoDB
	router *mux.Router
}

func NewServer(db data.EventsMongoDB) *RestServer {
	r := mux.NewRouter()

	server := &RestServer{
		db:     db,
		router: r,
	}

	r.HandleFunc("/health", healthCheck)
	r.HandleFunc("/event", server.createEvent).Methods(http.MethodPost)
	r.HandleFunc("/event", server.getEvent).Methods(http.MethodGet)
	r.HandleFunc("/event/all", server.listEvents).Methods(http.MethodGet)
	r.HandleFunc("/event", server.deleteEvent).Methods(http.MethodDelete)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := userFromBearer(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	return server
}
