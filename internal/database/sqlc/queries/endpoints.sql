-- name: ListApiEndpoint :many
SELECT * FROM api_endpoint
ORDER BY endpoint_name
LIMIT $1 OFFSET $2;

-- name: RegisterApiEndpoint :one 
INSERT INTO api_endpoint (endpoint_name, endpoint_description) 
VALUES ($1, $2)
RETURNING api_endpoint_id;

-- name: RegisterApiEndpoints :copyfrom
INSERT INTO api_endpoint (endpoint_name, endpoint_description) 
VALUES ($1, $2);

-- name: DeleteApiEndpointById :exec
DELETE FROM api_endpoint
WHERE api_endpoint_id = $1;