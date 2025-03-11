ALTER TABLE organization DROP FOREIGN KEY organization_ibfk_1;

ALTER TABLE organization_permission DROP FOREIGN KEY organization_permission_ibfk_1;

ALTER TABLE subscription DROP FOREIGN KEY subscription_ibfk_1;
ALTER TABLE subscription DROP FOREIGN KEY subscription_ibfk_2;

ALTER TABLE tier_base_pricing DROP FOREIGN KEY tier_base_pricing_ibfk_1;
ALTER TABLE tier_base_pricing DROP FOREIGN KEY tier_base_pricing_ibfk_2;

ALTER TABLE custom_endpoint_pricing DROP FOREIGN KEY custom_endpoint_pricing_ibfk_1;
ALTER TABLE custom_endpoint_pricing DROP FOREIGN KEY custom_endpoint_pricing_ibfk_2;

ALTER TABLE api_usage_summary DROP FOREIGN KEY api_usage_summary_ibfk_1;
ALTER TABLE api_usage_summary DROP FOREIGN KEY api_usage_summary_ibfk_2;
ALTER TABLE api_usage_summary DROP FOREIGN KEY api_usage_summary_ibfk_3;

ALTER TABLE billing_history DROP FOREIGN KEY billing_history_ibfk_1;

ALTER TABLE organization_permission DROP FOREIGN KEY organization_permission_ibfk_2;


DROP TABLE IF EXISTS api_usage_summary;
DROP TABLE IF EXISTS billing_history;
DROP TABLE IF EXISTS organization_permission;
DROP TABLE IF EXISTS custom_endpoint_pricing;
DROP TABLE IF EXISTS tier_base_pricing;
DROP TABLE IF EXISTS api_endpoint;
DROP TABLE IF EXISTS subscription;
DROP TABLE IF EXISTS organization;
DROP TABLE IF EXISTS subscription_tier;
DROP TABLE IF EXISTS organization_type;
DROP TABLE IF EXISTS resource_type;
