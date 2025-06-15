-- name: ListOrgType :many
SELECT * FROM organization_type
ORDER BY organization_type_name
LIMIT ? OFFSET ?;

-- name: CreateOrgType :execresult 
INSERT INTO organization_type (organization_type_name) 
VALUES (?);

-- name: CreateOrgTypes :copyfrom
INSERT INTO organization_type (organization_type_name) 
VALUES (?);

-- name: DeleteOrgTypeById :exec
DELETE FROM organization_type
WHERE organization_type_id = ?;