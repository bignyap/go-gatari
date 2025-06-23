package pricing

type CreateCustomPricingParams struct {
	CustomCostPerCall float64 `json:"custom_cost_per_call"`
	CustomRateLimit   int     `json:"custom_rate_limit"`
	SubscriptionID    int     `json:"subscription_id"`
	TierBasePricingID int     `json:"tier_base_pricing_id"`
}

type CreateCustomPricingOutput struct {
	ID int `json:"id"`
	CreateCustomPricingParams
}

type CreateTierPricingParams struct {
	BaseCostPerCall    float64 `json:"base_cost_per_call" validate:"required"`
	BaseRateLimit      *int    `json:"base_rate_limit"`
	ApiEndpointId      int     `json:"api_endpoint_id" validate:"required"`
	SubscriptionTierID int     `json:"subscription_tier_id" validate:"required"`
}

type CreateTierPricingOutput struct {
	ID int `json:"id"`
	CreateTierPricingParams
}

type CreateTierPricingWithTierName struct {
	CreateTierPricingOutput
	EndpointName string `json:"endpoint_name"`
}

type CreateTierPricingOutputWithCount struct {
	TotalItems int                             `json:"total_items"`
	Data       []CreateTierPricingWithTierName `json:"data"`
}
