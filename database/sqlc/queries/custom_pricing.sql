-- name: GetCustomPricing :many
SELECT * FROM custom_endpoint_pricing
WHERE subscription_id = ?
LIMIT ? OFFSET ?;

-- name: CreateCustomPricing :execresult 
INSERT INTO custom_endpoint_pricing (
    custom_cost_per_call, custom_rate_limit,
    subscription_id, tier_base_pricing_id
) 
VALUES (?, ?, ?, ?);

-- name: CreateCustomPricings :copyfrom
INSERT INTO custom_endpoint_pricing (
    custom_cost_per_call, custom_rate_limit,
    subscription_id, tier_base_pricing_id
) 
VALUES (?, ?, ?, ?);

-- name: DeleteCustomPricingById :exec
DELETE FROM custom_endpoint_pricing
WHERE custom_endpoint_pricing_id = ?;

-- name: DeleteCustomPricingBySubscriptionId :exec
DELETE FROM custom_endpoint_pricing
WHERE subscription_id = ?;