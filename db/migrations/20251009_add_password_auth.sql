-- migrate:up
ALTER TABLE users ADD COLUMN password_hash VARCHAR(255);
ALTER TABLE users ALTER COLUMN google_id DROP NOT NULL;

-- migrate:down
ALTER TABLE users ALTER COLUMN google_id SET NOT NULL;
ALTER TABLE users DROP COLUMN password_hash;
