-- name: ListOrgType :many
SELECT * FROM organization_type
ORDER BY organization_type_name
LIMIT $1 OFFSET $2;

-- name: CreateOrgType :one 
INSERT INTO organization_type (organization_type_name) 
VALUES ($1)
RETURNING organization_type_id;

-- name: CreateOrgTypes :copyfrom
INSERT INTO organization_type (organization_type_name) 
VALUES ($1);

-- name: DeleteOrgTypeById :exec
DELETE FROM organization_type
WHERE organization_type_id = $1;
