package scheduler

import (
	"log"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[scheduler] /test")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
