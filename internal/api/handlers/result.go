package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/service"
)

func (h *Handler) CheckResult(c *fiber.Ctx) error {
	jobID := c.Params("job_id")

	emails, err := service.FetchResult(c.Context(), jobID)
	if err == redis.Nil {
		return c.Type("html").SendString(`<p class="text-gray-500">⏳ Still waiting...</p>`)
	} else if err != nil {
		h.Log.WithError(err).Error("Failed to get result from Redis")
		return c.Status(500).SendString("Internal error")
	}

	html := "<ul class='text-green-600 space-y-1'>"
	for _, email := range emails {
		html += "<li>" + email + "</li>"
	}
	html += "</ul>"

	return c.Type("html").SendString(html)
}
