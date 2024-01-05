package operators

import (
	"context"
	"encoding/json"
	"fmt"
	"hookmq/config"
	"log"
	"time"

	cron "github.com/robfig/cron/v3"
)

func RunCron() {

	cronCtx := &config.ApplicationContext{
		Logger: &config.ServiceLogger{Service: "cron", ColorPrefix: config.ColorYellow},
	}

	cronCtx.Logger.Log("running cron")
	c := cron.New()
	c.AddFunc("* * * * *", func() { // Runs every minute
		currentTime := time.Now().Format(time.RFC3339)
		cronCtx.Logger.Log("\n")
		cronCtx.Logger.Log("---	Cron running	---")
		cronCtx.Logger.Log(fmt.Sprintf("Current time:\t%s", currentTime))
		jobs, err := GetDueJobsFromDynamoDB() // Retrieve jobs that are due for processing
		if err != nil {
			log.Println("Error fetching jobs:", err)
			return
		}

		// Publish to SQS
		for _, job := range jobs {
			// Serialize job data
			jobData, err := json.Marshal(job)
			if err != nil {
				cronCtx.Logger.Log(fmt.Sprintf("Error marshaling job: %v", err))
				continue // Skip this job and move to the next
			}

			// Prepare and send the message to SQS
			_, err = SQSClient.SendMessage(context.TODO(), string(jobData)) // Assuming SQSClient is an instance of SQSQueue
			if err != nil {
				cronCtx.Logger.Log(fmt.Sprintf("Failed to publish job to SQS: %v", err))
			} else {
				cronCtx.Logger.Log(fmt.Sprintf("Published job to SQS:\t%s", job.RowKey))
			}
		}

		cronCtx.Logger.Log("---	Cron stopped	---")
		cronCtx.Logger.Log("\n")
	})
	c.Start()
}
