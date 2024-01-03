package api

import (
	"log"
	"net/http"
)

// Job represents the structure of the job data
type Job struct {
	ID        string      `json:"id"`
	Payload   interface{} `json:"payload"` // Can be any JSON data
	URL       string      `json:"url"`
	ExecuteAt string      `json:"executeAt"`
	Status    string      `json:"Status"`
}

func JobHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("[api] /job")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job added successfully"))
}
