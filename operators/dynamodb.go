// operators/dynamodb.go
package operators

import (
	"context"
	"fmt"
	localConfig "hookmq/config"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
}

// ConnectDynamoDB initializes and sets up a DynamoDB client
func ConnectDynamoDB() {
	/*
		When using the AWS SDK for Go, if you have set the environment variables in (~/.aws/config)
		the SDK will automatically use these credentials. You don't need to manually specify them in your code.
	*/
	// Load the Shared AWS Configuration (~/.aws/config)
	// The profile name will be stored in an environment variable
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

	// Create DynamoDB client
	DynamoClient = dynamodb.NewFromConfig(cfg)

	// "Ping" DynamoDB to confirm connection
	// The table name to store processed message will be stored in an environment variable
	dynamoDBQueueTable := localConfig.ReadEnv("DYNAMODB_QUEUE_TABLE")
	if dynamoDBQueueTable == "" {
		log.Fatal("DYNAMODB_QUEUE_TABLE environment variable not set")
	}
	if err != nil {
		log.Fatalf("Failed to connect to DynamoDB: %v", err)
	} else {
		log.Print("Connected to DynamodDB successfully")
	}
}

func GetDueJobsFromDynamoDB() ([]Job, error) {
	// Define the current time as the threshold for jobs being due
	currentTime := time.Now().Format(time.RFC3339)

	// This uses Scan on the sort key which is inefficient
	// TODO: Update this query to use Global Secondary Index (GSI)
	// The table name to store processed message will be stored in an environment variable
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
		// DEBUG	-> Prints the raw item from the db table
		var job Job
		err = attributevalue.UnmarshalMap(item, &job)
		if err != nil {
			log.Printf("Failed to unmarshal DynamoDB item: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	// DEBUG 	-> This prints out the jobs that passed the conditional check
	// 			-> Use %+v for more detailed struct printing
	// log.Printf("Jobs: %+v", jobs)
	return jobs, nil
}
