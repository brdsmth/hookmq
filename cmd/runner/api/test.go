package runner

import (
	"log"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[runner] /test")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
