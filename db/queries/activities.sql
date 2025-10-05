-- name: GetActivitiesByUserAndDate :many
SELECT *
FROM activities
WHERE
    user_id = $1
    AND activity_date = $2
ORDER BY start_time ASC;

-- name: GetActivityByID :one
SELECT * FROM activities WHERE id = $1 LIMIT 1;

-- name: CreateActivity :one
INSERT INTO
    activities (
        user_id,
        title,
        start_time,
        end_time,
        category,
        roi,
        activity_date
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: UpdateActivity :one
UPDATE activities
SET
    title = $2,
    start_time = $3,
    end_time = $4,
    category = $5,
    roi = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteActivity :exec
DELETE FROM activities WHERE id = $1;

-- name: CheckActivityOverlap :one
-- Check if any existing activity overlaps with new time range
-- Binding: $1=user_id, $2=activity_date, $3=new_end_time, $4=new_start_time
-- Overlap: existing.start < new.end AND existing.end > new.start
SELECT EXISTS (
        SELECT 1
        FROM activities
        WHERE
            user_id = $1
            AND activity_date = $2
            AND start_time < $3::TIME
            AND end_time > $4::TIME
    ) AS has_overlap;

-- name: CheckActivityOverlapExcludingID :one
-- Binding: $1=user_id, $2=activity_date, $3=exclude_id, $4=new_end_time, $5=new_start_time
SELECT EXISTS (
        SELECT 1
        FROM activities
        WHERE
            user_id = $1
            AND activity_date = $2
            AND id <> $3
            AND start_time < $4::TIME
            AND end_time > $5::TIME
    ) AS has_overlap;