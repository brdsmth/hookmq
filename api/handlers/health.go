package handlers

import (
	"hookmq/config"
	"net/http"
)

// Now, your handlers can use the logger from the context
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	apiCtx := &config.ApplicationContext{
		Logger: &config.ServiceLogger{Service: "api", ColorPrefix: config.ColorGreen},
	}
	apiCtx.Logger.Log("--> /health")

	// Handle the request
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API OK"))
}
