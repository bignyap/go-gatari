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
