-- name: ListApiEndpoint :many
SELECT api_endpoint.*, resource_type.resource_type_name, permission_type.permission_code, permission_type.permission_name
FROM api_endpoint
INNER JOIN resource_type ON resource_type.resource_type_id = api_endpoint.resource_type_id
INNER JOIN permission_type ON permission_type.permission_code = api_endpoint.permission_code
ORDER BY api_endpoint_id DESC
LIMIT $1 OFFSET $2;

-- name: GetApiEndpointByName :one
SELECT api_endpoint.*, resource_type.resource_type_id, resource_type.resource_type_name, permission_type.permission_code, permission_type.permission_name
FROM api_endpoint
INNER JOIN resource_type ON resource_type.resource_type_id = api_endpoint.resource_type_id
INNER JOIN permission_type ON permission_type.permission_code = api_endpoint.permission_code
WHERE endpoint_name = $1;

-- name: GetApiEndpointById :one
SELECT api_endpoint.*, resource_type.resource_type_name, permission_type.permission_code, permission_type.permission_name
FROM api_endpoint
INNER JOIN resource_type ON resource_type.resource_type_id = api_endpoint.resource_type_id
INNER JOIN permission_type ON permission_type.permission_code = api_endpoint.permission_code
WHERE api_endpoint_id = $1;

-- name: ListApiEndpointsByResourceType :many
SELECT api_endpoint.*, resource_type.resource_type_name, permission_type.permission_code, permission_type.permission_name
FROM api_endpoint
INNER JOIN resource_type ON resource_type.resource_type_id = api_endpoint.resource_type_id
INNER JOIN permission_type ON permission_type.permission_code = api_endpoint.permission_code
WHERE api_endpoint.resource_type_id = $1
ORDER BY api_endpoint_id DESC;

-- name: RegisterApiEndpoint :one 
INSERT INTO api_endpoint (
  endpoint_name,
  endpoint_description,
  http_method,
  path_template,
  resource_type_id,
  permission_code
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING api_endpoint_id;

-- name: RegisterApiEndpoints :copyfrom
INSERT INTO api_endpoint (
  endpoint_name,
  endpoint_description,
  http_method,
  path_template,
  resource_type_id,
  permission_code
)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DeleteApiEndpointById :exec
DELETE FROM api_endpoint
WHERE api_endpoint_id = $1;

-- name: UpdateApiEndpointById :exec
UPDATE api_endpoint
SET
  endpoint_name = $2,
  endpoint_description = $3,
  http_method = $4,
  path_template = $5,
  resource_type_id = $6,
  permission_code = $7
WHERE api_endpoint_id = $1;

-- name: UpsertApiEndpointByName :one
INSERT INTO api_endpoint (
  endpoint_name,
  endpoint_description,
  http_method,
  path_template,
  resource_type_id,
  permission_code
)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (endpoint_name) DO UPDATE
SET
  endpoint_description = EXCLUDED.endpoint_description,
  http_method = EXCLUDED.http_method,
  path_template = EXCLUDED.path_template,
  resource_type_id = EXCLUDED.resource_type_id,
  permission_code = EXCLUDED.permission_code
RETURNING api_endpoint_id;

-- name: GetEndpointByName :one
SELECT
  api_endpoint_id AS id,
  endpoint_name AS name
FROM api_endpoint
WHERE endpoint_name = $1;