-- name: GetOrgPermission :many
SELECT * FROM organization_permission
WHERE organization_id = $1
LIMIT $2 OFFSET $3;

-- name: CreateOrgPermission :one 
INSERT INTO organization_permission (
    resource_type_id, permission_code, organization_id
) 
VALUES ($1, $2, $3)
RETURNING organization_id;

-- name: CreateOrgPermissions :copyfrom
INSERT INTO organization_permission (
    resource_type_id, permission_code, organization_id
) 
VALUES ($1, $2, $3);

-- name: DeleteOrgPermissionById :exec
DELETE FROM organization_permission
WHERE organization_permission_id = $1;

-- name: DeleteOrgPermissionByOrgId :exec
DELETE FROM organization_permission
WHERE organization_id = $1;

-- name: CheckOrgPermission :one
SELECT EXISTS (
  SELECT 1
  FROM organization_permission
  WHERE resource_type_id = $1
    AND permission_code = $2
);