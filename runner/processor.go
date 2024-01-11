package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hookmq/config"
	"hookmq/operators"
	"io"
	"net/http"
	"strings"
)

// sendHTTPRequest wraps http request for better error handling
func sendHTTPRequest(url string, payloadReader io.Reader) (string, error) {
	resp, err := http.Post(url, "application/json", payloadReader)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return resp.Status, nil
}

/*
The Processer handles each message from SQS and makes the POST reqeust defined in each job
*/
func Processor(job operators.Job) {
	runnerCtx := &config.ApplicationContext{
		Logger: &config.ServiceLogger{Service: "runner", ColorPrefix: config.ColorRed},
	}

	runnerCtx.Logger.Log(fmt.Sprintf("processing:\t\t\t%s", job.RowKey))

	var payloadReader io.Reader

	switch payload := job.Payload.(type) {
	case string:
		// If payload is already a JSON string, use it directly
		payloadReader = strings.NewReader(payload)
	default:
		// Else, marshal the payload to JSON
		payloadBytes, err := json.Marshal(job.Payload)
		if err != nil {
			runnerCtx.Logger.Log(fmt.Sprintf("Error marshaling payload for job %s: %v", job.JobID, err))
			return
		}
		payloadReader = bytes.NewReader(payloadBytes)
	}

	// Make an HTTP POST request with the JSON payload
	status, err := sendHTTPRequest(job.URL, payloadReader)
	if err != nil {
		runnerCtx.Logger.Log(fmt.Sprintf("http error %s: %v", job.JobID, err))
		operators.SQSClient.SendToDeadLetterQueue(context.TODO(), job)
		job.Status = "ERROR"
		operators.WriteToProcessed(job)
		return
	}

	runnerCtx.Logger.Log(fmt.Sprintf("sent:\t\t\t\t%s with status %s", job.JobID, status))

	// Update DynamoDB
	// Set the status of the job after the POST request is made
	job.Status = "PROCESSED"

	operators.WriteToProcessed(job)
}
