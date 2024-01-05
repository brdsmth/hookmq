package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"hookmq/config"
	"hookmq/operators"
	"log"
	"time"

	cron "github.com/robfig/cron/v3"
)

func RunCron() {

	schedulerCtx := &config.ApplicationContext{
		Logger: &config.ServiceLogger{Service: "scheduler", ColorPrefix: config.ColorPurple},
	}

	schedulerCtx.Logger.Log("running cron")
	c := cron.New()
	c.AddFunc("* * * * *", func() { // Runs every minute
		currentTime := time.Now().Format(time.RFC3339)
		schedulerCtx.Logger.Log("\n")
		schedulerCtx.Logger.Log("---	Cron running	---")
		schedulerCtx.Logger.Log(fmt.Sprintf("current time:\t\t%s", currentTime))
		jobs, err := operators.GetDueJobsFromDynamoDB() // Retrieve jobs that are due for processing
		if err != nil {
			log.Println("error fetching jobs:", err)
			return
		}

		// Publish to SQS
		for _, job := range jobs {
			// Serialize job data
			jobData, err := json.Marshal(job)
			if err != nil {
				schedulerCtx.Logger.Log(fmt.Sprintf("error marshaling job: %v", err))
				continue // Skip this job and move to the next
			}

			// Prepare and send the message to SQS
			_, err = operators.SQSClient.SendMessage(context.TODO(), string(jobData)) // Assuming SQSClient is an instance of SQSQueue
			if err != nil {
				schedulerCtx.Logger.Log(fmt.Sprintf("failed to publish job to SQS: %v", err))
			} else {
				schedulerCtx.Logger.Log(fmt.Sprintf("published to SQS:\t\t%s", job.RowKey))
			}
		}

		schedulerCtx.Logger.Log("---	Cron stopped	---")
		schedulerCtx.Logger.Log("\n")
	})
	c.Start()
}
