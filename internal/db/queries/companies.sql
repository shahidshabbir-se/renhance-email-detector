-- name: GetCompanyByDomain :one
SELECT * FROM companies
WHERE domain = $1;

-- name: CreateCompany :one
INSERT INTO companies (
    domain, organization, description, country, city
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;
