package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"hookmq/config"
	"hookmq/operators"
	"time"
)

func Listener(ctx context.Context) {
	runnerCtx := &config.ApplicationContext{
		Logger: &config.ServiceLogger{Service: "runner", ColorPrefix: config.ColorRed},
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Receive message
			result, err := operators.SQSClient.ReceiveMessage(context.TODO())

			if err != nil {
				runnerCtx.Logger.Log(fmt.Sprintf("error receiving message in sqs: %v", err))
				continue
			}

			// Process each message
			for _, message := range result.Messages {
				// Parse the message body
				var job operators.Job
				err := json.Unmarshal([]byte(*message.Body), &job)
				if err != nil {
					runnerCtx.Logger.Log(fmt.Sprintf("error parsing message in sqs: %v", err))
					continue // Skip to the next message
				}

				// Log the received message id
				runnerCtx.Logger.Log(fmt.Sprintf("message received in sqs:\t%s", job.RowKey))

				// Processs the job - make request to url defined in job
				Processor(job)

				// Delete the message after processing
				// Delete the message from the queue
				_, delErr := operators.SQSClient.DeleteMessage(ctx, message.ReceiptHandle)

				if delErr != nil {
					runnerCtx.Logger.Log(fmt.Sprintf("error deleting message from sqs: %v", delErr))
				}

				// Log the deleted message id
				runnerCtx.Logger.Log(fmt.Sprintf("message deleted from sqs:\t%s", job.RowKey))
			}
		}

		// Sleep for a while before next poll (if you want to rate limit the polling)
		time.Sleep(1 * time.Second)
	}
}
