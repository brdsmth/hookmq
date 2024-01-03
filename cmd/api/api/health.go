package api

import (
	"log"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[api] /health")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
