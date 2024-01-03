// cmd/api/main.go
package main

import (
	"fmt"
	"hookmq/cmd/api/api"
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

	fmt.Println("api has started...")
	fmt.Println("api is running on port 8081")

	http.ListenAndServe(":8081", nil)
}
