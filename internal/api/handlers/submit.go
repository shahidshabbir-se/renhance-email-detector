package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/service"
)

func (h *Handler) SubmitEmail(c *fiber.Ctx) error {
	company := c.FormValue("company")
	if company == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Company is required")
	}

	jobID := uuid.NewString()
	if err := service.EnqueueJob(c.Context(), jobID, company); err != nil {
		h.Log.WithError(err).Error("Failed to enqueue job")
		return c.Status(500).SendString("Internal error")
	}

	html := fmt.Sprintf(`
		<p class="text-blue-600">Job submitted with ID: <code>%s</code></p>
		<div hx-get="/result/%s" hx-trigger="every 5s" hx-swap="outerHTML"></div>
	`, jobID, jobID)

	return c.Type("html").SendString(html)
}
