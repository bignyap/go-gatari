-- +goose Up
ALTER TABLE api_endpoint ADD COLUMN access_type VARCHAR(20) NOT NULL DEFAULT 'paid' CHECK (access_type IN ('free', 'paid', 'private'));

-- +goose Down
ALTER TABLE api_endpoint DROP COLUMN IF EXISTS access_type;