package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/db"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Log *logrus.Logger
}

func New(log *logrus.Logger) *Handler {
	return &Handler{
		Log: log,
	}
}

func GetJobDetailsHandler(store db.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("id")
		results, err := store.GetJobDetails(c.Context(), jobID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(results)
	}
}
