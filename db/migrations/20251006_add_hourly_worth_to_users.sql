-- migrate:up
ALTER TABLE users ADD COLUMN hourly_worth DECIMAL(10,2) DEFAULT 25.00;

-- migrate:down
ALTER TABLE users DROP COLUMN hourly_worth;