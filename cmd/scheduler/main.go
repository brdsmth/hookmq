// cmd/scheduler/main.go
package main

import (
	"context"
	scheduler "hookmq/cmd/scheduler/api"
	"hookmq/internal/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Connect to SQS
	db.ConnectSQS()

	// Sending messaage to SQS client
	db.SQSClient.SendMessage(context.TODO(), "sending sending sending")

	// Initialize Mux router
	r := mux.NewRouter()

	r.HandleFunc("/health", scheduler.HealthHandler).Methods("GET")
	r.HandleFunc("/test", scheduler.TestHandler).Methods("GET")

	// Start the HTTP server for the API
	log.Println("scheduler listening on 8083")
	http.ListenAndServe(":8083", r)
}
