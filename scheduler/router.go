// scheduler/router.go
package scheduler

import (
	"hookmq/scheduler/handlers"

	"github.com/gorilla/mux"
)

func RegisterSchedulerRoutes(r *mux.Router) {
	apiRouter := r.PathPrefix("/scheduler").Subrouter()
	apiRouter.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
}
