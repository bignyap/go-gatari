-- name: ListApiEndpoint :many
SELECT api_endpoint.*, resource_type.resource_type_name
FROM api_endpoint
INNER JOIN resource_type ON resource_type.resource_type_id = api_endpoint.resource_type_id
ORDER BY api_endpoint_id DESC
LIMIT $1 OFFSET $2;

-- name: GetApiEndpointByName :one
SELECT api_endpoint.*, resource_type.resource_type_name
FROM api_endpoint
INNER JOIN resource_type ON resource_type.resource_type_id = api_endpoint.resource_type_id
WHERE endpoint_name = $1;

-- name: ListApiEndpointsByResourceType :many
SELECT api_endpoint.*, resource_type.resource_type_name
FROM api_endpoint
INNER JOIN resource_type ON resource_type.resource_type_id = api_endpoint.resource_type_id
WHERE api_endpoint.resource_type_id = $1
ORDER BY api_endpoint_id DESC;

-- name: RegisterApiEndpoint :one 
INSERT INTO api_endpoint (
  endpoint_name,
  endpoint_description,
  http_method,
  path_template,
  resource_type_id
)
VALUES ($1, $2, $3, $4, $5)
RETURNING api_endpoint_id;

-- name: RegisterApiEndpoints :copyfrom
INSERT INTO api_endpoint (
  endpoint_name,
  endpoint_description,
  http_method,
  path_template,
  resource_type_id
)
VALUES ($1, $2, $3, $4, $5);

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
  resource_type_id = $6
WHERE api_endpoint_id = $1;

-- name: UpsertApiEndpointByName :one
INSERT INTO api_endpoint (
  endpoint_name,
  endpoint_description,
  http_method,
  path_template,
  resource_type_id
)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (endpoint_name) DO UPDATE
SET
  endpoint_description = EXCLUDED.endpoint_description,
  http_method = EXCLUDED.http_method,
  path_template = EXCLUDED.path_template,
  resource_type_id = EXCLUDED.resource_type_id
RETURNING api_endpoint_id;