-- name: ListResourceType :many
SELECT * FROM resource_type
ORDER BY resource_type_name
LIMIT $1 OFFSET $2;

-- name: CreateResourceType :one 
INSERT INTO resource_type (
    resource_type_name, resource_type_code, resource_type_description
) 
VALUES ($1, $2, $3)
RETURNING resource_type_id;

-- name: CreateResourceTypes :copyfrom
INSERT INTO resource_type (
    resource_type_name, resource_type_code, resource_type_description
) 
VALUES ($1, $2, $3);

-- name: DeleteResourceTypeById :exec
DELETE FROM resource_type
WHERE resource_type_id = $1;