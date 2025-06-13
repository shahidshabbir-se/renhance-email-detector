package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db/sqlc"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/types"
)

func DetectEmails(domain string) ([]string, error) {
	req, _ := http.NewRequest("GET", "https://api.hunter.io/v2/domain-search", nil)
	q := req.URL.Query()
	q.Add("domain", domain)
	q.Add("api_key", os.Getenv("HUNTER_API_KEY"))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.HunterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	emails := make([]string, 0, len(result.Data.Emails))
	for _, e := range result.Data.Emails {
		emails = append(emails, e.Value)
	}
	return emails, nil
}

func SaveResult(ctx context.Context, q *sqlc.Queries, jobID, company, email string) error {
	emailText := pgtype.Text{String: email, Valid: email != ""}
	return q.InsertResult(ctx, sqlc.InsertResultParams{
		JobID:   jobID,
		Company: company,
		Email:   emailText,
	})
}

func CacheResult(ctx context.Context, redisClient *redis.Client, jobID string, emails []string) error {
	if len(emails) == 0 {
		return nil
	}
	key := fmt.Sprintf("result:%s", jobID)
	data, _ := json.Marshal(emails)
	return redisClient.Set(ctx, key, data, 24*time.Hour).Err()
}
