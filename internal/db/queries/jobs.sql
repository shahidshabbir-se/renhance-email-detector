-- name: CreateJob :one
INSERT INTO jobs (id, domain)
VALUES ($1, $2)
RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1;
