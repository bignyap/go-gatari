-- +goose Up

ALTER TABLE organization_permission
ADD CONSTRAINT unique_org_permission UNIQUE (resource_type_id, permission_code, organization_id);

-- +goose Down

ALTER TABLE organization_permission
DROP CONSTRAINT IF EXISTS unique_org_permission;