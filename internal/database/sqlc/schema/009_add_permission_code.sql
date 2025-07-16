-- +goose Up

ALTER TABLE api_endpoint ADD COLUMN permission_code VARCHAR(10) NOT NULL REFERENCES permission_type(permission_code);

-- +goose Down

ALTER TABLE api_endpoint DROP COLUMN IF EXISTS permission_code;