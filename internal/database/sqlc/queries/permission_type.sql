-- name: ListPermissionTypes :many
SELECT * FROM permission_type
ORDER BY permission_name
LIMIT $1 OFFSET $2;

-- name: CreatePermissionType :one
INSERT INTO permission_type (
    permission_code, permission_name, permission_description
)
VALUES ($1, $2, $3)
RETURNING permission_code;

-- name: CreatePermissionTypes :copyfrom
INSERT INTO permission_type (
    permission_code, permission_name, permission_description
)
VALUES ($1, $2, $3);

-- name: DeletePermissionTypeByCode :exec
DELETE FROM permission_type
WHERE permission_code = $1;