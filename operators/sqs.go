// operators/sqs.go
package operators

import (
	"context"
	"log"

	localConfig "hookmq/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSQueue represents a wrapper for the SQS client
type SQSQueue struct {
	Client   *sqs.Client
	QueueURL string
}

var SQSClient *SQSQueue

// New creates a new SQSQueue
func NewSQS(client *sqs.Client, queueURL string) *SQSQueue {
	return &SQSQueue{
		Client:   client,
		QueueURL: queueURL,
	}
}

// ConnectSQS initializes and sets up an SQS client
func ConnectSQS() {
	awsConfigProfile := localConfig.ReadEnv("AWS_CONFIG_PROFILE")
	if awsConfigProfile == "" {
		log.Fatal("AWS_CONFIG_PROFILE environment variable not set")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(awsConfigProfile),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)
	queueURL := localConfig.ReadEnv("SQS_URL")
	if queueURL == "" {
		log.Fatal("SQS_URL environment variable not set")
	}

	// Create an instance of SQSQueue
	SQSClient = NewSQS(client, queueURL)
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
		MaxNumberOfMessages: 10, // Adjust based on your needs
		WaitTimeSeconds:     10, // Adjust based on your needs
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
