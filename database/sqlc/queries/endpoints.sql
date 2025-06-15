-- name: ListApiEndpoint :many
SELECT * FROM api_endpoint
ORDER BY endpoint_name
LIMIT ? OFFSET ?;

-- name: RegisterApiEndpoint :execresult 
INSERT INTO api_endpoint (endpoint_name, endpoint_description) 
VALUES (?, ?);

-- name: RegisterApiEndpoints :copyfrom
INSERT INTO api_endpoint (endpoint_name, endpoint_description) 
VALUES (?, ?);

-- name: DeleteApiEndpointById :exec
DELETE FROM api_endpoint
WHERE api_endpoint_id = ?;