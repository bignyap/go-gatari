-- name: ListSubscription :many
SELECT 
    subscription.*, subscription_tier.tier_name, 
    COUNT(subscription.subscription_tier_id) OVER() AS total_items  
FROM subscription
INNER JOIN subscription_tier ON subscription.subscription_tier_id = subscription_tier.subscription_tier_id
ORDER BY subscription.subscription_tier_id DESC
LIMIT ? OFFSET ?;

-- name: GetSubscriptionById :one
SELECT 
    subscription.*, subscription_tier.tier_name  
FROM subscription
INNER JOIN subscription_tier ON subscription.subscription_tier_id = subscription_tier.subscription_tier_id
WHERE subscription.subscription_id = ?;

-- name: GetSubscriptionByOrgId :many
SELECT 
    subscription.*, subscription_tier.tier_name,
    COUNT(subscription.subscription_tier_id) OVER() AS total_items  
FROM subscription
INNER JOIN subscription_tier ON subscription.subscription_tier_id = subscription_tier.subscription_tier_id
WHERE subscription.organization_id = ?
LIMIT ? OFFSET ?;


-- name: CreateSubscription :execresult 
INSERT INTO subscription (
    subscription_name, subscription_type, subscription_created_date,
    subscription_updated_date, subscription_start_date, subscription_api_limit, 
    subscription_expiry_date, subscription_description, subscription_status, 
    organization_id, subscription_tier_id
) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateSubscriptions :copyfrom
INSERT INTO subscription (
    subscription_name, subscription_type, subscription_created_date,
    subscription_updated_date, subscription_start_date, subscription_api_limit, 
    subscription_expiry_date, subscription_description, subscription_status, 
    organization_id, subscription_tier_id
) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateSubscription :execresult
UPDATE subscription
SET 
    subscription_name = ?,
    subscription_start_date = ?,
    subscription_api_limit = ?,
    subscription_expiry_date = ?,
    subscription_description = ?,
    subscription_status = ?,
    organization_id = ?,
    subscription_tier_id = ?
WHERE subscription_id = ?;

-- name: DeleteSubscriptionByOrgId :exec
DELETE FROM subscription
WHERE organization_id = ?;

-- name: DeleteSubscriptionById :exec
DELETE FROM subscription
WHERE subscription_id = ?;