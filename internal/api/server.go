package api

import (
	"os"

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

	log.Infof("Starting Fiber API on :%s", port)
	return app.Listen(":" + port)
}
