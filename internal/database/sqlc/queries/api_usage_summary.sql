-- name: GetUsageSummary :many
WITH filtered_usage AS (
  SELECT organization_id, subscription_id, api_endpoint_id,
  SUM(total_calls) AS total_calls,
  SUM(total_cost) AS total_cost
  FROM api_usage_summary
  WHERE 
    (sqlc.narg('org_id')::int IS NULL OR organization_id = sqlc.narg('org_id')) AND
    (sqlc.narg('sub_id')::int IS NULL OR subscription_id = sqlc.narg('sub_id')) AND
    (sqlc.narg('endpoint_id')::int IS NULL OR api_endpoint_id = sqlc.narg('endpoint_id')) AND
    (sqlc.narg('start_date')::int IS NULL OR usage_start_date >= sqlc.narg('start_date')) AND
    (sqlc.narg('end_date')::int IS NULL OR usage_end_date <= sqlc.narg('end_date'))
  GROUP BY organization_id, subscription_id, api_endpoint_id
  LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset')
)
SELECT 
  f.organization_id,
  org.organization_name,
  f.subscription_id,
  sub.subscription_name,
  f.api_endpoint_id,
  ae.endpoint_name,
  f.total_calls,
  f.total_cost
FROM filtered_usage f
JOIN api_endpoint ae ON f.api_endpoint_id = ae.api_endpoint_id
JOIN subscription sub ON f.subscription_id = sub.subscription_id
JOIN organization org ON f.organization_id = org.organization_id;

-- name: GetUsageSummaryGroupedByDay :many
WITH filtered_usage AS (
  SELECT organization_id, subscription_id, api_endpoint_id,
  EXTRACT(YEAR FROM TO_TIMESTAMP(usage_start_date))::INT AS usage_year,
  EXTRACT(MONTH FROM TO_TIMESTAMP(usage_start_date))::INT AS usage_month,
  EXTRACT(DAY FROM TO_TIMESTAMP(usage_start_date))::INT AS usage_day,
  SUM(total_calls) AS total_calls,
  SUM(total_cost) AS total_cost
  FROM api_usage_summary
  WHERE 
    (sqlc.narg('org_id')::int IS NULL OR organization_id = sqlc.narg('org_id')) AND
    (sqlc.narg('sub_id')::int IS NULL OR subscription_id = sqlc.narg('sub_id')) AND
    (sqlc.narg('endpoint_id')::int IS NULL OR api_endpoint_id = sqlc.narg('endpoint_id')) AND
    (sqlc.narg('start_date')::int IS NULL OR usage_start_date >= sqlc.narg('start_date')) AND
    (sqlc.narg('end_date')::int IS NULL OR usage_end_date <= sqlc.narg('end_date'))
  GROUP BY usage_year, usage_month, usage_day, organization_id, subscription_id, api_endpoint_id
  LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset')
)
SELECT 
  f.organization_id,
  org.organization_name,
  f.subscription_id,
  sub.subscription_name,
  f.api_endpoint_id,
  ae.endpoint_name,
  f.usage_year,
  f.usage_month,
  f.usage_day,
  f.total_calls,
  f.total_cost
FROM filtered_usage f
JOIN api_endpoint ae ON f.api_endpoint_id = ae.api_endpoint_id
JOIN subscription sub ON f.subscription_id = sub.subscription_id
JOIN organization org ON f.organization_id = org.organization_id
ORDER BY usage_year DESC, usage_month DESC, usage_day DESC;

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

-- name: IncrementUsage :exec
INSERT INTO api_usage_summary (
  usage_start_date, usage_end_date, total_calls, total_cost,
  subscription_id, api_endpoint_id, organization_id
)
VALUES (
  FLOOR(EXTRACT(EPOCH FROM now())/60)*60,
  FLOOR(EXTRACT(EPOCH FROM now())/60)*60 + 59,
  1, 0.0, $1, $2, $3
)
ON CONFLICT (usage_start_date, usage_end_date, api_endpoint_id, organization_id)
DO UPDATE SET total_calls = api_usage_summary.total_calls + 1;
