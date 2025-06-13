package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/api/handlers"
	"github.com/shahidshabbir-se/renhance-email-detector/templates/pages"
	"github.com/sirupsen/logrus"
)

func SetupRouter(handler *handlers.Handler, log *logrus.Logger) *fiber.App {
	app := fiber.New()

	app.Static("/static", "./web/static")

	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html")
		return pages.Home().Render(c.Context(), c.Response().BodyWriter())
	})

	app.Post("/submit", handler.SubmitEmail)
	app.Get("/result/:job_id", handler.CheckResult)

	return app
}
