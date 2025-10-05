-- name: GetActivitiesByUserAndDate :many
SELECT * FROM activities
WHERE user_id = $1 AND activity_date = $2
ORDER BY start_time ASC;

-- name: GetActivityByID :one
SELECT * FROM activities WHERE id = $1 LIMIT 1;

-- name: CreateActivity :one
INSERT INTO activities (user_id, title, start_time, end_time, category, roi, activity_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateActivity :one
UPDATE activities
SET title = $2, start_time = $3, end_time = $4, category = $5, roi = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteActivity :exec
DELETE FROM activities WHERE id = $1;
