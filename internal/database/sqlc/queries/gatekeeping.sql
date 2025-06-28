-- name: GetOrganizationByName :one
SELECT
  organization_id AS id,
  organization_name AS name
FROM organization
WHERE organization_name = $1 AND organization_active = TRUE;

-- name: GetEndpointByName :one
SELECT
  api_endpoint_id AS id,
  endpoint_name AS name
FROM api_endpoint
WHERE endpoint_name = $1;

-- name: GetActiveSubscription :one
SELECT
  subscription_id AS id,
  organization_id,
  subscription_api_limit AS api_limit,
  subscription_expiry_date AS expiry_timestamp,
  subscription_status AS active
FROM subscription
WHERE organization_id = $1
  AND subscription_status = TRUE
  AND EXISTS (
    SELECT 1 FROM tier_base_pricing tbp
    WHERE tbp.subscription_tier_id = subscription.subscription_tier_id
      AND tbp.api_endpoint_id = $2
  );

-- name: GetPricing :one
SELECT
  COALESCE(cep.custom_cost_per_call, tbp.base_cost_per_call) AS cost_per_call
FROM subscription
JOIN tier_base_pricing tbp
  ON subscription.subscription_tier_id = tbp.subscription_tier_id
  AND tbp.api_endpoint_id = $2
LEFT JOIN custom_endpoint_pricing cep
  ON cep.subscription_id = subscription.subscription_id
  AND cep.tier_base_pricing_id = tbp.tier_base_pricing_id
WHERE subscription.subscription_id = $1;

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