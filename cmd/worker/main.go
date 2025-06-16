package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/db"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db/sqlc"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/logger"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/service"
	"github.com/shahidshabbir-se/renhance-email-detector/pkg/utils"
)

func main() {
	_ = godotenv.Load()

	log := logrus.New()
	logger.InitLogger(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize Redis
	if err := db.InitRedis(ctx, log); err != nil {
		log.WithError(err).Fatal("❌ Failed to initialize Redis")
	}

	// Initialize PostgreSQL
	pool := db.InitPostgres(ctx, log)
	queries := sqlc.New(pool)

	// Load Hunter API key
	apiKey := utils.GetEnv("HUNTER_API_KEY", "")
	if apiKey == "" {
		log.Fatal("❌ HUNTER_API_KEY is not set in environment")
	}

	// Init Hunter API client and service
	client := service.NewHunterClient(apiKey)
	hunterService := service.NewHunterService(queries, client, db.RedisClient)

	log.Info("✅ Email detection worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info("👋 Worker shutdown requested, exiting")
			return
		default:
			if err := hunterService.PollAndProcess(ctx, "email_detection_jobs"); err != nil {
				log.WithError(err).Warn("⚠️ PollAndProcess failed")
				time.Sleep(2 * time.Second)
			}
		}
	}
}
