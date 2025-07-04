-- name: GetTierLevelMonthlyUsage :many

SELECT
  t.subscription_tier_id,
  t.tier_name,
  DATE_TRUNC('month', TO_TIMESTAMP(a.usage_start_date)) AS usage_month,
  SUM(a.total_calls) AS total_calls,
  SUM(a.total_cost) AS total_revenue,
  COUNT(DISTINCT s.subscription_id) AS active_subscriptions
FROM
  api_usage_summary a
JOIN
  subscription s ON s.subscription_id = a.subscription_id
JOIN
  subscription_tier t ON t.subscription_tier_id = s.subscription_tier_id
WHERE
  TO_TIMESTAMP(a.usage_start_date) >= DATE_TRUNC('month', NOW() - INTERVAL '1 month')
  AND TO_TIMESTAMP(a.usage_start_date) < DATE_TRUNC('month', NOW())
GROUP BY
  t.subscription_tier_id, t.tier_name, usage_month
ORDER BY
  usage_month DESC, t.tier_name;


-- name: GetSubscriptionQuotaUsage :many
SELECT * FROM v_subscription_quota_usage WHERE quota_exceeded = true;

-- name: GetQuotaUsageBySubscriptionID :one
SELECT *
FROM v_subscription_quota_usage
WHERE subscription_id = $1;

-- name: EndpointUsagePerOrganization :many

SELECT
  o.organization_name,
  a.api_endpoint_id,
  e.endpoint_name,
  SUM(a.total_calls) AS monthly_calls
FROM
  api_usage_summary a
JOIN
  organization o ON o.organization_id = a.organization_id
JOIN
  api_endpoint e ON e.api_endpoint_id = a.api_endpoint_id
WHERE
  TO_TIMESTAMP(a.usage_start_date) >= DATE_TRUNC('month', NOW())
  AND TO_TIMESTAMP(a.usage_start_date) < DATE_TRUNC('month', NOW()) + INTERVAL '1 month'
GROUP BY
  o.organization_name, a.api_endpoint_id, e.endpoint_name
ORDER BY
  monthly_calls DESC;

  -- name: TotalCallsPerSubscription :many

  SELECT
  s.subscription_id,
  s.subscription_api_limit,
  DATE_TRUNC('month', TO_TIMESTAMP(a.usage_start_date)) AS month,
  SUM(a.total_calls) AS total_calls_in_month
FROM
  api_usage_summary a
JOIN
  subscription s ON s.subscription_id = a.subscription_id
WHERE
  TO_TIMESTAMP(a.usage_start_date) >= DATE_TRUNC('month', NOW()) AND
  TO_TIMESTAMP(a.usage_start_date) < DATE_TRUNC('month', NOW()) + INTERVAL '1 month'
GROUP BY
  s.subscription_id,
  s.subscription_api_limit,
  month;

  -- name: SubscriptionUsageExceedingLimit :many

  SELECT *
FROM (
  SELECT
    s.subscription_id,
    s.subscription_api_limit,
    SUM(a.total_calls) AS calls_made
  FROM
    api_usage_summary a
  JOIN
    subscription s ON s.subscription_id = a.subscription_id
  WHERE
    TO_TIMESTAMP(a.usage_start_date) >= DATE_TRUNC('month', NOW()) AND
    TO_TIMESTAMP(a.usage_start_date) < DATE_TRUNC('month', NOW()) + INTERVAL '1 month'
  GROUP BY
    s.subscription_id,
    s.subscription_api_limit
) AS usage_summary
WHERE
  subscription_api_limit IS NOT NULL AND calls_made > subscription_api_limit;



