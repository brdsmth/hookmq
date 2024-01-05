package data

import (
	"context"
	"hookmq/operators"
	"log"
)

func ProcessMessages() {
	ctx := context.TODO()

	// Connect to SQS (this sets up the SQSClient)
	operators.ConnectSQS()

	// Receive a message
	msgResult, err := operators.SQSClient.ReceiveMessage(ctx)
	if err != nil {
		log.Fatalf("Unable to receive messages: %v", err)
	}

	for _, message := range msgResult.Messages {
		// Process the message
		log.Printf("Message received: %s\n", *message.Body)

		// Delete the message from the queue
		_, delErr := operators.SQSClient.DeleteMessage(ctx, message.ReceiptHandle)
		if delErr != nil {
			log.Printf("Got an error while trying to delete message from queue: %v", delErr)
		} else {
			log.Println("Message deleted successfully")
		}
	}
}
