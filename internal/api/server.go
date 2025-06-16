package api

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/api/handlers"
	"github.com/sirupsen/logrus"
)

func StartServer(log *logrus.Logger) error {
	handler := handlers.New(log)
	app := SetupRouter(handler, log)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Infof("Starting Fiber API on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Fiber server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	time.Sleep(1 * time.Second) // optional short delay before shutdown

	if err := app.Shutdown(); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
		return err
	}

	log.Info("Server gracefully stopped.")
	return nil
}
