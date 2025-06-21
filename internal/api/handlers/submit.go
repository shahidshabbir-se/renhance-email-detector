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
	jobResultsURL := fmt.Sprintf("/result/%s", jobID)

	html := fmt.Sprintf(`
        <div
            style="padding: 16px; border-radius: 8px; text-align: center; max-width: 500px; margin: 20px auto 0px auto; background-color: transparent; color: #f0f0f0;"
        >
            <h2 style="color: #00FF94; font-size: 16px; margin-bottom: 15px; font-weight: bold;">Job Submitted Successfully!</h2>
            <p style="line-height: 1.6; font-size: 12px; margin-bottom: 20px;">
                Your request is being processed. You can track its progress at:
            </p>
            <p style="margin-top: 10px; margin-bottom: 10px;">
                <a
                    href="%s"
                    style="display: inline-block; padding: 8px 15px; background-color: #4a79ee; color: white; text-decoration: none; border-radius: 6px; font-weight: 500;"
                >
                    Go to Job Results
                </a>
            </p>
        </div>
    `, jobResultsURL)

	return c.Type("html").SendString(html)
}
