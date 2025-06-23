-- +goose Up

CREATE TABLE resource_type (
  resource_type_id SERIAL PRIMARY KEY,
  resource_type_code VARCHAR(10) UNIQUE NOT NULL,
  resource_type_name VARCHAR(50) NOT NULL,
  resource_type_description TEXT
);

CREATE TABLE organization_type (
  organization_type_id SERIAL PRIMARY KEY,
  organization_type_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE subscription_tier (
  subscription_tier_id SERIAL PRIMARY KEY,
  tier_name VARCHAR(50) NOT NULL,
  tier_archived BOOLEAN DEFAULT false NOT NULL,
  tier_description TEXT,
  tier_created_at INTEGER NOT NULL,
  tier_updated_at INTEGER NOT NULL
);

CREATE TABLE organization (
  organization_id SERIAL PRIMARY KEY,
  organization_name VARCHAR(100) UNIQUE NOT NULL,
  organization_created_at INTEGER NOT NULL,
  organization_updated_at INTEGER NOT NULL,
  organization_realm VARCHAR(100) NOT NULL,
  organization_country VARCHAR(50),
  organization_support_email VARCHAR(256) NOT NULL,
  organization_active BOOLEAN DEFAULT true,
  organization_report_q BOOLEAN DEFAULT false,
  organization_config TEXT,
  organization_type_id INTEGER NOT NULL REFERENCES organization_type(organization_type_id)
);

CREATE TABLE subscription (
  subscription_id SERIAL PRIMARY KEY,
  subscription_name VARCHAR(255) UNIQUE NOT NULL,
  subscription_type VARCHAR(255) NOT NULL,
  subscription_created_date INTEGER NOT NULL,
  subscription_updated_date INTEGER NOT NULL,
  subscription_start_date INTEGER NOT NULL,
  subscription_api_limit INTEGER,
  subscription_expiry_date INTEGER,
  subscription_description TEXT,
  subscription_status BOOLEAN DEFAULT true,
  organization_id INTEGER NOT NULL REFERENCES organization(organization_id),
  subscription_tier_id INTEGER NOT NULL REFERENCES subscription_tier(subscription_tier_id)
);

CREATE TABLE api_endpoint (
  api_endpoint_id SERIAL PRIMARY KEY,
  endpoint_name VARCHAR(255) UNIQUE NOT NULL,
  endpoint_description TEXT
);

CREATE TABLE tier_base_pricing (
  tier_base_pricing_id SERIAL PRIMARY KEY,
  base_cost_per_call FLOAT NOT NULL,
  base_rate_limit INTEGER,
  api_endpoint_id INTEGER NOT NULL REFERENCES api_endpoint(api_endpoint_id),
  subscription_tier_id INTEGER NOT NULL REFERENCES subscription_tier(subscription_tier_id),
  CONSTRAINT unique_api_tier UNIQUE (api_endpoint_id, subscription_tier_id)
);

CREATE TABLE custom_endpoint_pricing (
  custom_endpoint_pricing_id SERIAL PRIMARY KEY,
  custom_cost_per_call FLOAT NOT NULL,
  custom_rate_limit INTEGER NOT NULL,
  subscription_id INTEGER NOT NULL REFERENCES subscription(subscription_id),
  tier_base_pricing_id INTEGER NOT NULL REFERENCES tier_base_pricing(tier_base_pricing_id)
);

CREATE TABLE organization_permission (
  organization_permission_id SERIAL PRIMARY KEY,
  resource_type_id INTEGER NOT NULL REFERENCES resource_type(resource_type_id),
  permission_code VARCHAR(50) NOT NULL,
  organization_id INTEGER NOT NULL REFERENCES organization(organization_id)
);

CREATE TABLE billing_history (
  billing_id SERIAL PRIMARY KEY,
  billing_start_date INTEGER NOT NULL,
  billing_end_date INTEGER NOT NULL,
  total_amount_due FLOAT NOT NULL,
  total_calls INTEGER NOT NULL,
  payment_status VARCHAR(50) NOT NULL DEFAULT 'Pending',
  payment_date INTEGER,
  billing_created_at INTEGER NOT NULL,
  subscription_id INTEGER NOT NULL REFERENCES subscription(subscription_id)
);

CREATE TABLE api_usage_summary (
  usage_summary_id SERIAL PRIMARY KEY,
  usage_start_date INTEGER NOT NULL,
  usage_end_date INTEGER NOT NULL,
  total_calls INTEGER NOT NULL,
  total_cost FLOAT NOT NULL,
  subscription_id INTEGER NOT NULL REFERENCES subscription(subscription_id),
  api_endpoint_id INTEGER NOT NULL REFERENCES api_endpoint(api_endpoint_id),
  organization_id INTEGER NOT NULL REFERENCES organization(organization_id)
  -- ,CONSTRAINT unique_org_endpoint_period UNIQUE (usage_start_date, usage_end_date, api_endpoint_id, organization_id)
);

-- +goose Down

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