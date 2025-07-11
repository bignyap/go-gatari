-- +goose Up
-- Add foreign key constraint to organization_permission table
ALTER TABLE organization_permission
  ADD CONSTRAINT fk_permission_code
  FOREIGN KEY (permission_code)
  REFERENCES permission_type(permission_code);

-- +goose Down
-- Remove foreign key constraint first
ALTER TABLE organization_permission
  DROP CONSTRAINT IF EXISTS fk_permission_code;
