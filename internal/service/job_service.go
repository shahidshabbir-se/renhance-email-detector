package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/db"
)

func EnqueueJob(ctx context.Context, jobID, company string) error {
	job := map[string]string{
		"job_id":       jobID,
		"company_name": company,
		"created_at":   time.Now().String(),
	}
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return db.RedisClient.RPush(ctx, "email_detection_jobs", data).Err()
}

func FetchResult(ctx context.Context, jobID string) ([]string, error) {
	val, err := db.RedisClient.Get(ctx, "result:"+jobID).Result()
	if err != nil {
		return nil, err
	}

	var emails []string
	if err := json.Unmarshal([]byte(val), &emails); err != nil {
		return nil, fmt.Errorf("invalid result format: %w", err)
	}

	return emails, nil
}
