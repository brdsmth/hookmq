package scheduler

import (
	"log"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[scheduler] /health")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
