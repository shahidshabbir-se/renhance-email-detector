package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/db/sqlc"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/types"
)

type HunterService struct {
	Store  *sqlc.Queries
	Client *HunterClient
	Redis  *redis.Client
}

func NewHunterService(store *sqlc.Queries, client *HunterClient, redisClient *redis.Client) *HunterService {
	return &HunterService{
		Store:  store,
		Client: client,
		Redis:  redisClient,
	}
}

func (hs *HunterService) PollAndProcess(ctx context.Context, queue string) error {
	log.Infof("🔄 Waiting for job on Redis queue: %s", queue)

	result, err := hs.Redis.BLPop(ctx, 0, queue).Result()
	if err != nil {
		log.Errorf("❌ Redis BLPop failed: %v", err)
		return err
	}

	if len(result) < 2 {
		log.Error("⚠️ Invalid job format in Redis (less than 2 elements)")
		return errors.New("invalid Redis job format")
	}

	rawJob := result[1]
	log.Infof("📦 Raw job payload: %s", rawJob)

	var job types.Job
	if err := json.Unmarshal([]byte(rawJob), &job); err != nil {
		log.Errorf("❌ Failed to unmarshal job: %v", err)
		return err
	}

	log.Infof("✅ Parsed job - Company: %s | JobID: %s", job.Company, job.ID)

	resp, err := hs.Client.SearchDomain(job.Company)
	if err != nil {
		log.Errorf("❌ Hunter.io API error for %s: %v", job.Company, err)
		return err
	}

	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	log.Infof("📨 Hunter.io response for %s:\n%s", job.Company, string(respJSON))

	return hs.persistHunterData(ctx, job.ID, job.Company, resp)
}

func (hs *HunterService) persistHunterData(ctx context.Context, jobID string, domain string, data *types.HunterResponse) error {
	if data == nil || data.Data.Domain == "" || len(data.Data.Emails) == 0 {
		log.Warnf("Invalid or empty Hunter.io response for job ID: %s", jobID)
		return nil
	}

	company, err := hs.Store.GetCompanyByDomain(ctx, domain)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			company, err = hs.Store.CreateCompany(ctx, sqlc.CreateCompanyParams{
				Domain:       domain,
				Organization: pgText(data.Data.Organization),
				Description:  pgText(data.Data.Description),
				Country:      pgText(data.Data.Country),
				City:         pgText(data.Data.City),
			})
			if err != nil {
				log.Errorf("Failed to create company: %v", err)
				return err
			}
		} else {
			log.Errorf("Failed to query company: %v", err)
			return err
		}
	}

	_, err = hs.Store.CreateJob(ctx, sqlc.CreateJobParams{
		ID:     jobID,
		Domain: domain,
	})
	if err != nil {
		log.Errorf("Failed to create job: %v", err)
		return err
	}

	for _, email := range data.Data.Emails {
		_, err := hs.Store.CreateContact(ctx, sqlc.CreateContactParams{
			CompanyID:   company.ID,
			Email:       email.Value,
			FirstName:   pgText(email.FirstName),
			LastName:    pgText(email.LastName),
			Position:    pgText(email.Position),
			Department:  pgText(email.Department),
			LinkedinUrl: pgText(email.LinkedIn),
		})
		if err != nil {
			log.Warnf("Failed to insert contact: %v", err)
		}
	}

	return hs.Store.LinkJobResult(ctx, sqlc.LinkJobResultParams{
		JobID:     jobID,
		CompanyID: company.ID,
	})
}

func pgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}
