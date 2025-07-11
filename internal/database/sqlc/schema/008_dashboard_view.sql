-- +goose Up
CREATE VIEW dashboard_summary_view AS
SELECT
  (SELECT COUNT(*) FROM resource_type) AS resource_type_count,
  (SELECT COUNT(*) FROM api_endpoint) AS api_endpoint_count,
  (SELECT COUNT(*) FROM organization) AS organization_count,
  (SELECT COUNT(*) FROM organization WHERE organization_active = true) AS active_organization_count,
  (SELECT COUNT(*) FROM subscription_tier) AS subscription_tier_count,
  (SELECT COUNT(*) FROM subscription_tier WHERE tier_archived = false) AS active_subscription_tier_count,
  (SELECT COUNT(*) FROM subscription) AS subscription_count,
  (SELECT COUNT(*) FROM subscription WHERE subscription_status = true) AS active_subscription_count,
  (SELECT COUNT(*) FROM billing_history) AS billing_history_count,
  (SELECT COUNT(*) FROM api_usage_summary) AS usage_summary_count,
  (SELECT COUNT(*) FROM organization_permission) AS organization_permission_count,
  (SELECT COUNT(*) FROM permission_type) AS permission_type_count;

-- +goose Down
DROP VIEW IF EXISTS dashboard_summary_view;