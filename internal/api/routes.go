package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/api/handlers"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db"
	"github.com/shahidshabbir-se/renhance-email-detector/pkg/utils"
	"github.com/shahidshabbir-se/renhance-email-detector/templates/pages"
	"github.com/sirupsen/logrus"
)

func SetupRouter(handler *handlers.Handler, store db.Store, log *logrus.Logger) *fiber.App {
	adminPassword := utils.GetEnv("ADMIN_PASSWORD", "password123")
	app := fiber.New()

	app.Static("/static", "./web/static")

	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html")
		return pages.Home().Render(c.Context(), c.Response().BodyWriter())
	})

	app.Post("/submit", handler.SubmitEmail)
	app.Get("/result/:job_id", handlers.GetJobDetailsPage(store))
	app.Get("/metrics", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": adminPassword,
		},
	}), monitor.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	return app
}
