-- +goose Up

ALTER TABLE tier_base_pricing
ADD COLUMN cost_mode VARCHAR(10) NOT NULL DEFAULT 'fixed';

ALTER TABLE tier_base_pricing
ADD CONSTRAINT tier_base_pricing_cost_mode_check CHECK (cost_mode IN ('fixed', 'dynamic'));

ALTER TABLE custom_endpoint_pricing
ADD COLUMN cost_mode VARCHAR(10) NOT NULL DEFAULT 'fixed';

ALTER TABLE custom_endpoint_pricing
ADD CONSTRAINT custom_endpoint_pricing_cost_mode_check CHECK (cost_mode IN ('fixed', 'dynamic'));

-- +goose Down

ALTER TABLE tier_base_pricing
DROP CONSTRAINT IF EXISTS tier_base_pricing_cost_mode_check;

ALTER TABLE tier_base_pricing
DROP COLUMN IF EXISTS cost_mode;

ALTER TABLE custom_endpoint_pricing
DROP CONSTRAINT IF EXISTS custom_endpoint_pricing_cost_mode_check;

ALTER TABLE custom_endpoint_pricing
DROP COLUMN IF EXISTS cost_mode;