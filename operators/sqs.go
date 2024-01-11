// operators/sqs.go
package operators

import (
	"context"
	"encoding/json"
	"log"

	localConfig "hookmq/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSQueue represents a wrapper for the SQS client
type SQSQueue struct {
	Client             *sqs.Client
	QueueURL           string
	DeadLetterQueueURL string
}

var SQSClient *SQSQueue

// New creates a new SQSQueue
func NewSQS(client *sqs.Client, queueURL string, deadLetterQueueUrl string) *SQSQueue {
	return &SQSQueue{
		Client:             client,
		QueueURL:           queueURL,
		DeadLetterQueueURL: deadLetterQueueUrl,
	}
}

// ConnectSQS initializes and sets up an SQS client
func ConnectSQS() {
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

	client := sqs.NewFromConfig(cfg)
	queueURL := localConfig.ReadEnv("SQS_URL")
	if queueURL == "" {
		log.Fatal("SQS_URL environment variable not set")
	}

	deadLetterQueueURL := localConfig.ReadEnv("SQS_DL_URL")
	if queueURL == "" {
		log.Fatal("SQS_DL_URL environment variable not set")
	}

	// Create an instance of SQSQueue
	SQSClient = NewSQS(client, queueURL, deadLetterQueueURL)
}

// SendMessage sends a message to the SQS queue
func (q *SQSQueue) SendMessage(ctx context.Context, message string) (*sqs.SendMessageOutput, error) {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(q.QueueURL),
	}
	return q.Client.SendMessage(ctx, input)
}

// ReceiveMessage receives messages from the SQS queue
func (q *SQSQueue) ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.QueueURL),
		MaxNumberOfMessages: int32(*aws.Int64(10)),
		WaitTimeSeconds:     int32(*aws.Int64(20)), // Enable long polling
	}
	return q.Client.ReceiveMessage(ctx, input)
}

// DeleteMessage deletes a message from the SQS queue
func (q *SQSQueue) DeleteMessage(ctx context.Context, receiptHandle *string) (*sqs.DeleteMessageOutput, error) {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.QueueURL),
		ReceiptHandle: receiptHandle,
	}
	return q.Client.DeleteMessage(ctx, input)
}

// SendToDeadLetterQueue sends a message to the Dead Letter Queue defined in env
func (q *SQSQueue) SendToDeadLetterQueue(ctx context.Context, job Job) (*sqs.SendMessageOutput, error) {

	// Serialize the job or create a message
	message, err := json.Marshal(job)
	if err != nil {
		log.Printf("Failed to marshal job for dead letter queue: %v", err)
		return nil, err
	}

	// Convert message to a string before using with aws.String
	messageStr := string(message)

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(messageStr),
		QueueUrl:    aws.String(q.DeadLetterQueueURL),
	}

	return q.Client.SendMessage(ctx, input)
}
