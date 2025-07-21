-- name: GetTierPricingByTierId :many
SELECT 
    tier_base_pricing.*, api_endpoint.endpoint_name,
    COUNT(tier_base_pricing_id) OVER() AS total_items
FROM tier_base_pricing
INNER JOIN api_endpoint ON tier_base_pricing.api_endpoint_id = api_endpoint.api_endpoint_id
WHERE subscription_tier_id = $1
LIMIT $2 OFFSET $3;

-- name: CreateTierPricing :one 
INSERT INTO tier_base_pricing (subscription_tier_id, api_endpoint_id, base_cost_per_call, base_rate_limit, cost_mode) 
VALUES ($1, $2, $3, $4, $5)
RETURNING tier_base_pricing_id;

-- name: CreateTierPricings :copyfrom
INSERT INTO tier_base_pricing (subscription_tier_id, api_endpoint_id, base_cost_per_call, base_rate_limit, cost_mode) 
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateTierPricingByTierId :execresult
UPDATE tier_base_pricing
SET 
    base_cost_per_call = $1,
    base_rate_limit = $2,
    api_endpoint_id = $3,
    cost_mode = $4
WHERE subscription_tier_id = $5;

-- name: UpdateTierPricingById :execresult
UPDATE tier_base_pricing
SET 
    base_cost_per_call = $1,
    base_rate_limit = $2,
    api_endpoint_id = $3,
    cost_mode = $4
WHERE tier_base_pricing_id = $5;

-- name: DeleteTierPricingByTierId :one
DELETE FROM tier_base_pricing
WHERE subscription_tier_id = $1
RETURNING subscription_tier_id;

-- name: DeleteTierPricingById :one
DELETE FROM tier_base_pricing
WHERE tier_base_pricing_id = $1
RETURNING subscription_tier_id;

-- name: GetPricing :one
SELECT
  COALESCE(cep.custom_cost_per_call, tbp.base_cost_per_call, 0)::double precision AS cost_per_call,
  COALESCE(cep.cost_mode, tbp.cost_mode, 'fixed') AS cost_mode
FROM subscription
JOIN tier_base_pricing tbp
  ON subscription.subscription_tier_id = tbp.subscription_tier_id
  AND tbp.api_endpoint_id = $2
LEFT JOIN custom_endpoint_pricing cep
  ON cep.subscription_id = subscription.subscription_id
  AND cep.tier_base_pricing_id = tbp.tier_base_pricing_id
WHERE subscription.subscription_id = $1;