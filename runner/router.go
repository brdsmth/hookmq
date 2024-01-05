// runner/router.go
package runner

import (
	"hookmq/runner/handlers"

	"github.com/gorilla/mux"
)

func RegisterRunnerRoutes(r *mux.Router) {
	apiRouter := r.PathPrefix("/runner").Subrouter()
	apiRouter.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
}
