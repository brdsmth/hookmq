// cmd/api/main.go
package main

import (
	"fmt"
	"hookmq/cmd/api/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize Mux router
	r := mux.NewRouter()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is hookmq api")
	})

	r.HandleFunc("/health", api.HealthHandler).Methods("GET")
	r.HandleFunc("/job", api.JobHandler).Methods("GET")

	// Start the HTTP server for the API
	log.Println("api listening on 8081")
	http.ListenAndServe(":8081", r)
}
