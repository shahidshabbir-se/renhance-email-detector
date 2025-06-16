-- name: LinkJobResult :exec
INSERT INTO job_results (job_id, company_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;
