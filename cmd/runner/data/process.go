package data

import (
	"context"
	"hookmq/internal/db"
	"log"
)

func ProcessMessages() {
	ctx := context.TODO()

	// Connect to SQS (this sets up the SQSClient)
	db.ConnectSQS()

	// Receive a message
	msgResult, err := db.SQSClient.ReceiveMessage(ctx)
	if err != nil {
		log.Fatalf("Unable to receive messages: %v", err)
	}

	for _, message := range msgResult.Messages {
		// Process the message
		log.Printf("Message received: %s\n", *message.Body)

		// Delete the message from the queue
		_, delErr := db.SQSClient.DeleteMessage(ctx, message.ReceiptHandle)
		if delErr != nil {
			log.Printf("Got an error while trying to delete message from queue: %v", delErr)
		} else {
			log.Println("Message deleted successfully")
		}
	}
}
