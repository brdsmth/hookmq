// api/router.go
package api

import (
	"hookmq/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterApiRoutes(r *mux.Router) {
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	apiRouter.HandleFunc("/queue", handlers.QueueHandler).Methods("POST")
	apiRouter.HandleFunc("/test", handlers.TestHandler).Methods("POST")
}
