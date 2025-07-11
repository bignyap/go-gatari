-- +goose Up
ALTER TABLE api_endpoint ADD COLUMN http_method VARCHAR(10) NOT NULL DEFAULT 'GET';
ALTER TABLE api_endpoint ADD COLUMN path_template TEXT NOT NULL DEFAULT '/';
ALTER TABLE api_endpoint ADD COLUMN resource_type_id INTEGER NOT NULL REFERENCES resource_type(resource_type_id);

-- +goose Down
ALTER TABLE api_endpoint DROP COLUMN IF EXISTS http_method;
ALTER TABLE api_endpoint DROP COLUMN IF EXISTS path_template;
ALTER TABLE api_endpoint DROP COLUMN IF EXISTS resource_type_id;

