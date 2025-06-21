package handlers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/shahidshabbir-se/renhance-email-detector/internal/db/sqlc"
	"github.com/shahidshabbir-se/renhance-email-detector/internal/types"
	"github.com/shahidshabbir-se/renhance-email-detector/templates/pages"
)

func GetJobDetailsPage(store sqlc.Querier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("job_id")
		if jobID == "" {
			return fiber.NewError(http.StatusBadRequest, "missing job ID")
		}

		rows, err := store.GetJobDetails(context.Background(), jobID)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "failed to fetch job details")
		}
		if len(rows) == 0 {
			return fiber.NewError(http.StatusNotFound, "job not found")
		}

		company := types.Company{
			ID:           rows[0].CompanyID,
			Domain:       rows[0].Domain,
			Organization: pgTextToString(rows[0].Organization),
			Description:  pgTextToString(rows[0].Description),
			Country:      pgTextToString(rows[0].Country),
			City:         pgTextToString(rows[0].City),
		}

		contacts := make([]types.Contact, 0)
		for _, row := range rows {
			if row.ContactID.Valid {
				contacts = append(contacts, types.Contact{
					ID:         int(row.ContactID.Int32),
					Email:      pgTextToString(row.Email),
					FirstName:  pgTextToPtr(row.FirstName),
					LastName:   pgTextToPtr(row.LastName),
					Position:   pgTextToPtr(row.Position),
					Department: pgTextToPtr(row.Department),
					LinkedIn:   pgTextToPtr(row.LinkedinUrl),
				})
			}
		}

		job := types.JobDetails{
			JobID:     rows[0].JobID,
			JobDomain: rows[0].JobDomain,
			Company:   company,
			Contacts:  contacts,
		}

		// Render the Templ page with the full job details struct
		c.Type("html")
		return pages.Results(job).Render(c.Context(), c.Response().BodyWriter())
	}
}

// Helpers

func pgTextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

func pgTextToPtr(t pgtype.Text) *string {
	if t.Valid {
		return &t.String
	}
	return nil
}
