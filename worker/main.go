package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/repository"
)

var (
	ctx           = context.Background()
	redisClient   *redis.Client
	pgPool        *pgxpool.Pool
	dbQueries     *repository.Queries
	mailgunAPIKey string
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

	mailgunAPIKey = os.Getenv("MAILGUN_API_KEY")
	if mailgunAPIKey == "" {
		log.Fatal("MAILGUN_API_KEY environment variable is required")
	}
}

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})

	// Initialize PostgreSQL connection pool
	var err error
	pgPool, err = pgxpool.New(ctx, getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/emaildb"))
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	dbQueries = repository.New(pgPool)

	log.Println("Worker started, waiting for jobs...")

	for {
		jobJSON, err := redisClient.BLPop(ctx, 0*time.Second, "email_detection_jobs").Result()
		if err != nil {
			log.Printf("Error fetching job from Redis: %v", err)
			continue
		}

		if len(jobJSON) < 2 {
			continue
		}

		var job Job
		if err := json.Unmarshal([]byte(jobJSON[1]), &job); err != nil {
			log.Printf("Failed to parse job JSON: %v", err)
			continue
		}

		log.Printf("Processing job: %s for company: %s", job.JobID, job.CompanyName)

		emails, err := detectEmails(job.CompanyName)
		if err != nil {
			log.Printf("Error detecting emails: %v", err)
			continue
		}

		firstEmail := ""
		if len(emails) > 0 {
			firstEmail = emails[0]
		}

		err = saveResult(job.JobID, job.CompanyName, firstEmail)
		if err != nil {
			log.Printf("Error saving result: %v", err)
		}

		resultKey := fmt.Sprintf("result:%s", job.JobID)
		emailsJSON, _ := json.Marshal(emails)
		err = redisClient.Set(ctx, resultKey, emailsJSON, 24*time.Hour).Err()
		if err != nil {
			log.Printf("Error setting result in Redis: %v", err)
		}

		log.Printf("Job %s processed with %d emails found", job.JobID, len(emails))
	}
}

func detectEmails(company string) ([]string, error) {
	// TODO: Integrate with actual Mailgun API here using mailgunAPIKey
	// Dummy response for now:
	return []string{
		fmt.Sprintf("contact@%s.com", company),
		fmt.Sprintf("info@%s.com", company),
	}, nil
}

func saveResult(jobID, company, email string) error {
	emailText := pgtype.Text{
		String: email,
		Valid:  email != "",
	}

	return dbQueries.InsertResult(ctx, repository.InsertResultParams{
		JobID:   jobID,
		Company: company,
		Email:   emailText,
	})
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
