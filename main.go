// main.go
package main

import (
	"context"
	"fmt"
	"hookmq/api"
	"hookmq/config"
	"hookmq/operators"
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
		Logger: &config.ServiceLogger{Service: "hookmq", ColorPrefix: config.ColorCyan},
	}

	// Connect to SQS
	operators.ConnectSQS()

	// Connect to DynamoDB
	operators.ConnectDynamoDB()

	// Initiate cron
	scheduler.RunCron()

	// Start the SQS listener in a goroutine
	go func() {
		runner.Listener(context.Background())
	}()

	// Initialize Mux router
	r := mux.NewRouter()

	// Register API routes with the logger from the 'api' package
	api.RegisterApiRoutes(r)

	// Start the HTTP server for the publisher microservice
	port := config.ReadEnv("PORT")
	if port == "" {
		port = "8081" // default
	}

	hookmqCtx.Logger.Log(fmt.Sprintf("listening on port %s", port))
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start hookmq on %s: %v\n", port, err)
	}
}
