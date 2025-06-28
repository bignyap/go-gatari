-- name: ListSubscription :many
SELECT 
    subscription.*, subscription_tier.tier_name, 
    COUNT(subscription.subscription_tier_id) OVER() AS total_items  
FROM subscription
INNER JOIN subscription_tier ON subscription.subscription_tier_id = subscription_tier.subscription_tier_id
ORDER BY subscription.subscription_tier_id DESC
LIMIT $1 OFFSET $2;

-- name: GetSubscriptionById :one
SELECT 
    subscription.*, subscription_tier.tier_name  
FROM subscription
INNER JOIN subscription_tier ON subscription.subscription_tier_id = subscription_tier.subscription_tier_id
WHERE subscription.subscription_id = $1;

-- name: GetSubscriptionByOrgId :many
SELECT 
    subscription.*, subscription_tier.tier_name,
    COUNT(subscription.subscription_tier_id) OVER() AS total_items  
FROM subscription
INNER JOIN subscription_tier ON subscription.subscription_tier_id = subscription_tier.subscription_tier_id
WHERE subscription.organization_id = $1
LIMIT $2 OFFSET $3;

-- name: CreateSubscription :one 
INSERT INTO subscription (
    subscription_name, subscription_type, subscription_created_date,
    subscription_updated_date, subscription_start_date, subscription_api_limit, 
    subscription_expiry_date, subscription_description, subscription_status, 
    organization_id, subscription_tier_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING subscription_id;

-- name: CreateSubscriptions :copyfrom
INSERT INTO subscription (
    subscription_name, subscription_type, subscription_created_date,
    subscription_updated_date, subscription_start_date, subscription_api_limit, 
    subscription_expiry_date, subscription_description, subscription_status, 
    organization_id, subscription_tier_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: UpdateSubscription :execresult
UPDATE subscription
SET 
    subscription_name = $1,
    subscription_start_date = $2,
    subscription_api_limit = $3,
    subscription_expiry_date = $4,
    subscription_description = $5,
    subscription_status = $6,
    organization_id = $7,
    subscription_tier_id = $8
WHERE subscription_id = $9;

-- name: DeleteSubscriptionByOrgId :exec
DELETE FROM subscription
WHERE organization_id = $1;

-- name: DeleteSubscriptionById :exec
DELETE FROM subscription
WHERE subscription_id = $1;