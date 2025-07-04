-- name: GetBillingHistoryByOrgId :many
SELECT * FROM billing_history
WHERE subscription_id IN (
    SELECT subscription_id FROM subscription
    WHERE organization_id = $1
)
LIMIT $2 OFFSET $3;

-- name: GetBillingHistoryBySubId :many
SELECT * FROM billing_history
WHERE subscription_id = $1
LIMIT $2 OFFSET $3;

-- name: GetBillingHistoryById :many
SELECT * FROM billing_history
WHERE billing_id = $1
LIMIT $2 OFFSET $3;

-- name: CreateBillingHistory :one 
INSERT INTO billing_history (
    billing_start_date, billing_end_date, total_amount_due,
    total_calls, payment_status, payment_date, 
    billing_created_at, subscription_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING billing_id;

-- name: CreateBillingHistories :copyfrom
INSERT INTO billing_history (
    billing_start_date, billing_end_date, total_amount_due,
    total_calls, payment_status, payment_date, 
    billing_created_at, subscription_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);


-- INSERT monthly billing records for subscriptions that had usage last month
-- name: InsertMonthlyBillingRecords :exec
INSERT INTO billing_history (
  billing_start_date,
  billing_end_date,
  total_amount_due,
  total_calls,
  billing_created_at,
  subscription_id
)
SELECT
  EXTRACT(EPOCH FROM date_trunc('month', NOW() - INTERVAL '1 month'))::INTEGER AS billing_start_date,
  EXTRACT(EPOCH FROM date_trunc('month', NOW()))::INTEGER AS billing_end_date,
  COALESCE(SUM(a.total_cost), 0) AS total_amount_due,
  COALESCE(SUM(a.total_calls), 0) AS total_calls,
  EXTRACT(EPOCH FROM NOW())::INTEGER AS billing_created_at,
  a.subscription_id
FROM
  api_usage_summary a
WHERE
  TO_TIMESTAMP(a.usage_start_date) >= date_trunc('month', NOW() - INTERVAL '1 month')
  AND TO_TIMESTAMP(a.usage_start_date) < date_trunc('month', NOW())
GROUP BY
  a.subscription_id
ON CONFLICT DO NOTHING;

