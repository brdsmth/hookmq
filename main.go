// main.go
package main

import (
	"hookmq/api"
	"hookmq/config"
	"hookmq/runner"
	"hookmq/scheduler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorReset  = "\033[0m"
)

func main() {
	hookmqCtx := &config.ApplicationContext{
		Logger: &config.ServiceLogger{Service: "hookmq", ColorPrefix: config.ColorYellow},
	}

	// // Connect to SQS
	// db.ConnectSQS()

	// // Sending messaage to SQS client
	// db.SQSClient.SendMessage(context.TODO(), "sending sending sending")

	// Initialize Mux router
	r := mux.NewRouter()

	// Register API routes with the logger from the 'api' package
	api.RegisterApiRoutes(r)
	runner.RegisterRunnerRoutes(r)
	scheduler.RegisterSchedulerRoutes(r)

	hookmqCtx.Logger.Log("listening on 8081")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatalf("Failed to start hookmq on 8081: %v\n", err)
	}
}
