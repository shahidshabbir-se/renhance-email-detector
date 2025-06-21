package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/api"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/logger"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()

	log := logrus.New()
	logger.InitLogger(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := db.InitRedis(ctx, log); err != nil {
		log.WithError(err).Fatal("Failed to initialize Redis")
	}

	pool := db.InitPostgres(ctx, log)
	store := db.NewStore(pool)

	if err := api.StartServer(log, store); err != nil {
		log.WithError(err).Fatal("API server failed")
	}
}
