-- name: InsertHealthcheck :one
INSERT INTO healthchecks DEFAULT VALUES
RETURNING id, created_at;

-- name: ListHealthchecks :many
SELECT id, created_at
FROM healthchecks
ORDER BY id DESC
LIMIT $1;

