-- +goose Up
ALTER TABLE subscription ADD COLUMN subscription_quota_reset_interval VARCHAR(20)
    CHECK (subscription_quota_reset_interval IN ('monthly', 'yearly', 'total')) DEFAULT 'total';
ALTER TABLE subscription ADD COLUMN subscription_billing_model VARCHAR(20)
    CHECK (subscription_billing_model IN ('flat', 'usage', 'hybrid')) DEFAULT 'usage';
ALTER TABLE subscription ADD COLUMN subscription_billing_interval VARCHAR(20)
    CHECK (subscription_billing_interval IN ('monthly', 'yearly', 'once')) DEFAULT 'once';

-- +goose Down
ALTER TABLE subscription DROP COLUMN IF EXISTS subscription_quota_reset_interval;
ALTER TABLE subscription DROP COLUMN IF EXISTS subscription_billing_model;
ALTER TABLE subscription DROP COLUMN IF EXISTS subscription_billing_interval;
