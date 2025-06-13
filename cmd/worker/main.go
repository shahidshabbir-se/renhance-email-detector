package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db/sqlc"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/tasks"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/types"
)

var (
	ctx         = context.Background()
	log         = logrus.New()
	redisClient *redis.Client
	dbQueries   *sqlc.Queries
)

func main() {
	_ = godotenv.Load()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	pgPool, _ := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	defer pgPool.Close()
	dbQueries = sqlc.New(pgPool)

	log.Info("Worker started")

	for {
		jobJSON, err := redisClient.BLPop(ctx, 0*time.Second, "email_detection_jobs").Result()
		if err != nil || len(jobJSON) < 2 {
			continue
		}

		var job types.Job
		if err := json.Unmarshal([]byte(jobJSON[1]), &job); err != nil {
			log.WithError(err).Error("Invalid job JSON")
			continue
		}

		emails, err := tasks.DetectEmails(job.CompanyName)
		if err != nil {
			log.WithError(err).Error("Email detection failed")
			continue
		}

		first := ""
		if len(emails) > 0 {
			first = emails[0]
		}

		if err := tasks.SaveResult(ctx, dbQueries, job.JobID, job.CompanyName, first); err != nil {
			log.WithError(err).Error("Failed to save result")
		}
		if err := tasks.CacheResult(ctx, redisClient, job.JobID, emails); err != nil {
			log.WithError(err).Error("Failed to cache result")
		}

		log.WithFields(logrus.Fields{
			"job_id": job.JobID,
			"emails": len(emails),
		}).Info("Job completed")
	}
}
