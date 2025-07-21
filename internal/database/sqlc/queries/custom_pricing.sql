-- name: GetCustomPricing :many
SELECT * FROM custom_endpoint_pricing
WHERE subscription_id = $1
LIMIT $2 OFFSET $3;

-- name: CreateCustomPricing :one 
INSERT INTO custom_endpoint_pricing (
    custom_cost_per_call, custom_rate_limit,
    subscription_id, tier_base_pricing_id, cost_mode
) 
VALUES ($1, $2, $3, $4, $5)
RETURNING custom_endpoint_pricing_id;

-- name: CreateCustomPricings :copyfrom
INSERT INTO custom_endpoint_pricing (
    custom_cost_per_call, custom_rate_limit,
    subscription_id, tier_base_pricing_id, cost_mode
) 
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteCustomPricingById :one
DELETE FROM custom_endpoint_pricing
WHERE custom_endpoint_pricing_id = $1
RETURNING subscription_id;

-- name: DeleteCustomPricingBySubscriptionId :one
DELETE FROM custom_endpoint_pricing
WHERE subscription_id = $1
RETURNING subscription_id;