package handlers

import (
	"hookmq/config"
	"net/http"
)

// Now, your handlers can use the logger from the context
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	serviceCtx := &config.ApplicationContext{
		Logger: &config.ServiceLogger{Service: "scheduler", ColorPrefix: config.ColorPurple},
	}

	serviceCtx.Logger.Log("hello")
	// Handle the request
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SCHEDULER OK"))
}
