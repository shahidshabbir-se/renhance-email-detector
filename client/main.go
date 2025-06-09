package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

type Job struct {
	JobID       string `json:"job_id"`
	CompanyName string `json:"company_name"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on environment variables")
	}
}

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}()

	// Example: add jobs to the queue
	jobs := []Job{
		{JobID: "job-1", CompanyName: "example"},
		{JobID: "job-2", CompanyName: "openai"},
	}

	for _, job := range jobs {
		err := enqueueJob(job)
		if err != nil {
			log.Printf("Failed to enqueue job %s: %v", job.JobID, err)
		} else {
			log.Printf("Enqueued job %s for company %s", job.JobID, job.CompanyName)
		}
	}
}

func enqueueJob(job Job) error {
	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = redisClient.RPush(ctx, "email_detection_jobs", jobJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to push job to Redis: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
