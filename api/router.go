// api/router.go
package api

import (
	"hookmq/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterApiRoutes(r *mux.Router) {
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
}
