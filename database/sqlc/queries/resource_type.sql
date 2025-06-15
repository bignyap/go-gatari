-- name: ListResourceType :many
SELECT * FROM resource_type
ORDER BY resource_type_name
LIMIT ? OFFSET ?;

-- name: CreateResourceType :execresult 
INSERT INTO resource_type (resource_type_name, resource_type_code, resource_type_description) 
VALUES (?, ?, ?);

-- name: CreateResourceTypes :copyfrom
INSERT INTO resource_type (resource_type_name, resource_type_code, resource_type_description) 
VALUES (?, ?, ?);

-- name: DeleteResourceTypeById :exec
DELETE FROM resource_type
WHERE resource_type_id = ?;