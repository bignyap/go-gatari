-- +goose Up
CREATE OR REPLACE VIEW v_subscription_quota_usage AS
SELECT
  s.subscription_id,
  s.subscription_name,
  s.subscription_api_limit,
  s.subscription_quota_reset_interval,
  s.subscription_billing_model,
  s.subscription_billing_interval,
  COALESCE(SUM(a.total_calls), 0)::INT AS calls_used,
  (s.subscription_api_limit - COALESCE(SUM(a.total_calls), 0)::INT) AS calls_remaining,
  CASE
    WHEN s.subscription_api_limit IS NULL THEN NULL
    WHEN COALESCE(SUM(a.total_calls), 0)::INT >= s.subscription_api_limit THEN true
    ELSE false
  END AS quota_exceeded
FROM
  subscription s
LEFT JOIN api_usage_summary a
  ON s.subscription_id = a.subscription_id
  AND (
    (s.subscription_quota_reset_interval = 'monthly' AND TO_TIMESTAMP(a.usage_start_date) >= DATE_TRUNC('month', NOW()))
    OR
    (s.subscription_quota_reset_interval = 'yearly' AND TO_TIMESTAMP(a.usage_start_date) >= DATE_TRUNC('year', NOW()))
    OR
    (s.subscription_quota_reset_interval = 'total')
  )
WHERE
  s.subscription_status = true
  AND (s.subscription_expiry_date IS NULL OR s.subscription_expiry_date > EXTRACT(EPOCH FROM NOW()))
GROUP BY
  s.subscription_id,
  s.subscription_name,
  s.subscription_api_limit,
  s.subscription_quota_reset_interval,
  s.subscription_billing_model,
  s.subscription_billing_interval;

-- +goose Down
DROP VIEW IF EXISTS v_subscription_quota_usage;