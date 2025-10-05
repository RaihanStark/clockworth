-- migrate:up
-- Composite index to optimize overlap queries
-- This speeds up CheckActivityOverlap and CheckActivityOverlapExcludingID queries
CREATE INDEX idx_activities_overlap ON activities(user_id, activity_date, start_time, end_time);

-- migrate:down
DROP INDEX IF EXISTS idx_activities_overlap;
