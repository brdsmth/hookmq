// operators/dynamodb.go
package operators

import (
	"context"
	"encoding/json"
	"fmt"
	localConfig "hookmq/config"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DynamoClient *dynamodb.Client

type Job struct {
	ExecuteAt string      `dynamodbav:"ExecuteAt"`
	JobID     string      `dynamodbav:"JobID"`
	Payload   interface{} `dynamodbav:"Payload"` // Flexible for any JSON structure
	RowKey    string      `dynamodbav:"RowKey"`
	URL       string      `dynamodbav:"URL"`
	Status    string      `json:"Status"`
}

// ConnectDynamoDB initializes and sets up a DynamoDB client
func ConnectDynamoDB() {
	hookmqCtx := &localConfig.ApplicationContext{
		Logger: &localConfig.ServiceLogger{Service: "hookmq", ColorPrefix: localConfig.ColorCyan},
	}

	// Retrieve AWS credentials from environment variables
	awsAccessKeyID := localConfig.ReadEnv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := localConfig.ReadEnv("AWS_SECRET_ACCESS_KEY")
	awsDefaultRegion := localConfig.ReadEnv("AWS_DEFAULT_REGION")

	// Check if the necessary environment variables are set
	if awsAccessKeyID == "" || awsSecretAccessKey == "" {
		log.Fatal("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables must be set")
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
		config.WithRegion(awsDefaultRegion),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create DynamoDB client
	DynamoClient = dynamodb.NewFromConfig(cfg)

	// "Ping" DynamoDB to confirm connection
	dynamoDBQueueTable := localConfig.ReadEnv("DYNAMODB_QUEUE_TABLE")
	if dynamoDBQueueTable == "" {
		log.Fatal("DYNAMODB_QUEUE_TABLE environment variable not set")
	}
	if err != nil {
		log.Fatalf("Failed to connect to DynamoDB: %v", err)
	} else {
		hookmqCtx.Logger.Log("connected to dynamodb")
	}
}

func GetDueJobsFromDynamoDB() ([]Job, error) {
	// Define the current time as the threshold for jobs being due
	currentTime := time.Now().Format(time.RFC3339)

	// TODO: Update this query to use Global Secondary Index (GSI)
	// This uses Scan on the sort key which is inefficient
	dynamoDBQueueTable := localConfig.ReadEnv("DYNAMODB_QUEUE_TABLE")
	if dynamoDBQueueTable == "" {
		log.Fatal("DYNAMODB_QUEUE_TABLE environment variable not set")
	}
	tableName := dynamoDBQueueTable
	result, err := DynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("ExecuteAt <= :currentTime"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":currentTime": &types.AttributeValueMemberS{Value: currentTime},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query DynamoDB: %w", err)
	}

	// Decode the result items into a slice of Jobs
	var jobs []Job
	for _, item := range result.Items {
		var job Job
		err = attributevalue.UnmarshalMap(item, &job)
		if err != nil {
			log.Printf("Failed to unmarshal DynamoDB item: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func WriteToProcessed(job Job) {
	hookmqCtx := &localConfig.ApplicationContext{
		Logger: &localConfig.ServiceLogger{Service: "hookmq", ColorPrefix: localConfig.ColorCyan},
	}

	// Marshal the Payload to a JSON string
	postProcessPayloadBytes, err := json.Marshal(job.Payload)
	if err != nil {
		log.Printf("error processing job:\t%s", job.RowKey)
		return
	}
	postProcessPayloadString := string(postProcessPayloadBytes)
	currentTime := time.Now().Format(time.RFC3339)

	item := map[string]types.AttributeValue{
		"RowKey":     &types.AttributeValueMemberS{Value: job.RowKey},
		"JobID":      &types.AttributeValueMemberS{Value: job.JobID},
		"Payload":    &types.AttributeValueMemberS{Value: postProcessPayloadString},
		"URL":        &types.AttributeValueMemberS{Value: job.URL},
		"Status":     &types.AttributeValueMemberS{Value: job.Status},
		"ExecutedAt": &types.AttributeValueMemberS{Value: currentTime},
	}

	// Write the updated job to the `processed` Table
	dynamoDBProcessedTable := localConfig.ReadEnv("DYNAMODB_PROCESSED_TABLE")
	if dynamoDBProcessedTable == "" {
		log.Fatal("DYNAMODB_PROCESSED_TABLE environment variable not set")
	}
	processedTableName := dynamoDBProcessedTable
	_, err = DynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(processedTableName),
		Item:      item,
	})
	if err != nil {
		hookmqCtx.Logger.Log(fmt.Sprintf("failed to write to processed dynamodb: %v", err))
		return
	}

	hookmqCtx.Logger.Log(fmt.Sprintf("processed job:\t\t\t%s", job.JobID))

	DeleteFromQueue(job)
}

func DeleteFromQueue(job Job) {

	hookmqCtx := &localConfig.ApplicationContext{
		Logger: &localConfig.ServiceLogger{Service: "hookmq", ColorPrefix: localConfig.ColorCyan},
	}

	// Define the table from which to delete the item
	dynamoDBQueueTable := localConfig.ReadEnv("DYNAMODB_QUEUE_TABLE")
	if dynamoDBQueueTable == "" {
		log.Fatal("DYNAMODB_QUEUE_TABLE environment variable not set")
	}

	// Define the key of the item to be deleted
	key := map[string]types.AttributeValue{
		"RowKey":    &types.AttributeValueMemberS{Value: job.RowKey},
		"ExecuteAt": &types.AttributeValueMemberS{Value: job.ExecuteAt},
	}

	// Perform the deletion
	_, err := DynamoClient.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(dynamoDBQueueTable),
		Key:       key,
	})
	if err != nil {
		hookmqCtx.Logger.Log(fmt.Sprintf("failed to delete item from %s: %v", dynamoDBQueueTable, err))
	} else {
		hookmqCtx.Logger.Log(fmt.Sprintf("successfully deleted from db\t%s from %s", job.RowKey, dynamoDBQueueTable))
	}
}
