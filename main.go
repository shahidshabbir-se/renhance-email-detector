package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/api"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/logger"
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

	if err := api.StartServer(log); err != nil {
		log.WithError(err).Fatal("API server failed")
	}
}
