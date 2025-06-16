-- name: CreateContact :one
INSERT INTO contacts (
    company_id, email, first_name, last_name, position, department, linkedin_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ListContactsByCompanyID :many
SELECT * FROM contacts
WHERE company_id = $1
ORDER BY created_at DESC;
