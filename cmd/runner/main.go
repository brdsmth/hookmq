// cmd/runner/main.go
package main

import (
	runner "hookmq/cmd/runner/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize Mux router
	r := mux.NewRouter()

	r.HandleFunc("/health", runner.HealthHandler).Methods("GET")
	r.HandleFunc("/test", runner.TestHandler).Methods("GET")

	// Connect to SQS
	// db.ConnectSQS()

	// message, err := db.SQSClient.ReceiveMessage(context.TODO())
	// if err != nil {
	// 	log.Printf("cannot read message: %s", err)
	// }
	// log.Print(message)

	// TODO: Messages are received in batch mode
	// output, err := db.SQSClient.DeleteMessage(context.TODO(), message.Messages[0].ReceiptHandle)
	// if err != nil {
	// 	log.Printf("cannot read message: %s", err)
	// }
	// log.Print(output)

	// data.ProcessMessages()

	// Start the HTTP server for the Runner (testing)
	log.Println("runner listening on 8082")
	http.ListenAndServe(":8082", r)
}
