package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	localConfig "hookmq/config"
	"hookmq/operators"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gonanoid "github.com/matoous/go-nanoid"
)

// Job represents the structure of the job data
type Job struct {
	ID        string      `json:"id"`
	Payload   interface{} `json:"payload"` // Can be any JSON data
	URL       string      `json:"url"`
	ExecuteAt string      `json:"executeAt"`
	Status    string      `json:"Status"`
}

func Queue(w http.ResponseWriter, r *http.Request) {
	apiCtx := &localConfig.ApplicationContext{
		Logger: &localConfig.ServiceLogger{Service: "api", ColorPrefix: localConfig.ColorGreen},
	}
	apiCtx.Logger.Log("--> /queue")

	var job Job
	var err error
	// Identify job
	id, err := gonanoid.Nanoid(6)
	if err != nil {
		apiCtx.Logger.Log(fmt.Sprintf("error adding nanoid: %s", err))
		return
	}
	job.ID = id
	apiCtx.Logger.Log(fmt.Sprintf("add job:\t%s", id))

	// Decode the job from the request body
	err = json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		apiCtx.Logger.Log(fmt.Sprintf("cancel add job [decode]:\t%s", id))
		return
	}

	// Check if URL is provided
	if job.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		apiCtx.Logger.Log(fmt.Sprintf("cancel add job [empty]:\t%s", id))
		return
	}

	// Check if Payload is provided
	if job.Payload == nil {
		http.Error(w, "Payload is required", http.StatusBadRequest)
		return
	}

	// Marshal the Payload to a JSON string
	payloadBytes, err := json.Marshal(job.Payload)
	if err != nil {
		http.Error(w, "Error marshaling payload", http.StatusInternalServerError)
		apiCtx.Logger.Log(fmt.Sprintf("cancel add job [marshal]:\t%s", id))
		return
	}
	payloadString := string(payloadBytes)

	// Set ExecutionTime to now if not provided
	if job.ExecuteAt == "" {
		job.ExecuteAt = time.Now().Format(time.RFC3339)
	}

	// Set job status
	job.Status = "QUEUED"

	// Update DynamoDB
	partitionKey := fmt.Sprintf("queue::%s::%s", job.ID, job.ExecuteAt)

	// Prepare the item to write to DynamoDB
	item := map[string]types.AttributeValue{
		"RowKey":    &types.AttributeValueMemberS{Value: partitionKey},
		"JobID":     &types.AttributeValueMemberS{Value: job.ID},
		"Payload":   &types.AttributeValueMemberS{Value: payloadString}, // Serialized JSON string
		"URL":       &types.AttributeValueMemberS{Value: job.URL},
		"ExecuteAt": &types.AttributeValueMemberS{Value: job.ExecuteAt},
		"Status":    &types.AttributeValueMemberS{Value: job.Status},
	}

	// Write to DynamoDB
	// The table name to store processed message will be stored in an environment variable
	dynamoDBQueueTable := localConfig.ReadEnv("DYNAMODB_QUEUE_TABLE")
	if dynamoDBQueueTable == "" {
		log.Fatal("DYNAMODB_QUEUE_TABLE environment variable not set")
	}
	tableName := dynamoDBQueueTable
	_, err = operators.DynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		apiCtx.Logger.Log(fmt.Sprintf("failed to write to db: %v", err))
		http.Error(w, "Failed to add job", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job added successfully"))
}
