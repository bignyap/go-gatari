package pricing

type CreateCustomPricingParams struct {
	CustomCostPerCall float64 `json:"custom_cost_per_call" form:"custom_cost_per_call"`
	CustomRateLimit   int     `json:"custom_rate_limit" form:"custom_rate_limit"`
	SubscriptionID    int     `json:"subscription_id" form:"subscription_id"`
	TierBasePricingID int     `json:"tier_base_pricing_id" form:"tier_base_pricing_id"`
	CostMode          string  `json:"cost_mode" form:"cost_mode" validate:"required"`
}

type CreateCustomPricingOutput struct {
	ID int `json:"id"`
	CreateCustomPricingParams
}

type CreateTierPricingParams struct {
	BaseCostPerCall    float64 `json:"base_cost_per_call" form:"base_cost_per_call" validate:"required"`
	BaseRateLimit      *int    `json:"base_rate_limit" form:"base_rate_limit"`
	ApiEndpointId      int     `json:"api_endpoint_id" form:"api_endpoint_id" validate:"required"`
	SubscriptionTierID int     `json:"subscription_tier_id" form:"subscription_tier_id" validate:"required"`
	CostMode           string  `json:"cost_mode" form:"cost_mode" validate:"required"`
}

type CreateTierPricingOutput struct {
	ID int `json:"id"`
	CreateTierPricingParams
}

type CreateTierPricingWithTierName struct {
	CreateTierPricingOutput
	EndpointName string `json:"endpoint_name" form:"endpoint_name"`
}

type CreateTierPricingOutputWithCount struct {
	TotalItems int                             `json:"total_items"`
	Data       []CreateTierPricingWithTierName `json:"data"`
}
