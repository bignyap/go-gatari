-- name: GetApiUsageSummaryByOrgId :many
SELECT * FROM api_usage_summary
WHERE subscription_id IN (
    SELECT subscription_id FROM subscription s
    WHERE s.organization_id = $1
)
LIMIT $2 OFFSET $3;

-- name: GetApiUsageSummaryBySubId :many
SELECT * FROM api_usage_summary
WHERE subscription_id = $1
LIMIT $2 OFFSET $3;

-- name: GetApiUsageSummaryByEndpointId :many
SELECT * FROM api_usage_summary
WHERE api_endpoint_id = $1
LIMIT $2 OFFSET $3;

-- name: CreateApiUsageSummary :one 
INSERT INTO api_usage_summary (
    usage_start_date, usage_end_date, total_calls,
    total_cost, subscription_id, api_endpoint_id, 
    organization_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING usage_summary_id;

-- name: CreateApiUsageSummaries :copyfrom
INSERT INTO api_usage_summary (
    usage_start_date, usage_end_date, total_calls,
    total_cost, subscription_id, api_endpoint_id, 
    organization_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7);
