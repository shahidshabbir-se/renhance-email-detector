-- name: InsertResult :exec
INSERT INTO results (job_id, company, email, fetched_at) VALUES ($1, $2, $3, NOW());

-- name: GetResultByJobID :one
SELECT job_id, company, email, fetched_at FROM results WHERE job_id = $1;
