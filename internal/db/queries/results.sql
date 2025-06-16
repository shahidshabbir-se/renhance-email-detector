-- name: GetJobDetails :many
SELECT
    j.id AS job_id,
    j.domain AS job_domain,
    c.id AS company_id,
    c.domain,
    c.organization,
    c.description,
    c.country,
    c.city,
    ct.id AS contact_id,
    ct.email,
    ct.first_name,
    ct.last_name,
    ct.position,
    ct.department,
    ct.linkedin_url
FROM jobs j
JOIN job_results jr ON jr.job_id = j.id
JOIN companies c ON c.id = jr.company_id
LEFT JOIN contacts ct ON ct.company_id = c.id
WHERE j.id = $1;
