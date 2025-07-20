package pricing

import (
	"fmt"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func (h *PricingService) CreateTierPricingJSONValidation(c *gin.Context) ([]sqlcgen.CreateTierPricingsParams, error) {
	var inputs []CreateTierPricingParams

	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	outputs := make([]sqlcgen.CreateTierPricingsParams, 0, len(inputs))

	for i, input := range inputs {
		if err := h.Validator.Struct(input); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}

		output := sqlcgen.CreateTierPricingsParams{
			BaseCostPerCall:    input.BaseCostPerCall,
			BaseRateLimit:      converter.ToPgInt4(input.BaseRateLimit),
			ApiEndpointID:      int32(input.ApiEndpointId),
			SubscriptionTierID: int32(input.SubscriptionTierID),
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}

type TierPricingForm struct {
	ApiEndpointID      int32   `form:"api_endpoint_id" binding:"required"`
	SubscriptionTierID int32   `form:"subscription_tier_id" binding:"required"`
	BaseRateLimit      int32   `form:"base_rate_limit" binding:"required"`
	BaseCostPerCall    float64 `form:"base_cost_per_call" binding:"required"`
	CostMode           string  `form:"cost_mode" binding:"omitempty,oneof=fixed dynamic"`
}

func (h *PricingService) CreateTierPricingFormValidator(c *gin.Context) (*sqlcgen.CreateTierPricingParams, error) {
	var form TierPricingForm

	if err := c.ShouldBind(&form); err != nil {
		return nil, err
	}

	if form.CostMode == "" {
		form.CostMode = "fixed"
	}

	input := sqlcgen.CreateTierPricingParams{
		BaseCostPerCall:    form.BaseCostPerCall,
		SubscriptionTierID: form.SubscriptionTierID,
		BaseRateLimit:      converter.ToPgInt4(ptrInt(form.BaseRateLimit)),
		ApiEndpointID:      form.ApiEndpointID,
		CostMode:           form.CostMode,
	}

	return &input, nil
}

func ptrInt(v int32) *int {
	i := int(v)
	return &i
}

type CustomPricingForm struct {
	TierBasePricingID int32   `form:"tier_base_pricing_id" binding:"required"`
	SubscriptionID    int32   `form:"subscription_id" binding:"required"`
	CustomRateLimit   int32   `form:"custom_rate_limit" binding:"required"`
	CustomCostPerCall float64 `form:"custom_cost_per_call" binding:"required"`
	CostMode          string  `form:"cost_mode" binding:"omitempty,oneof=fixed dynamic"`
}

func (h *PricingService) CreateCustomPricingFormValidator(c *gin.Context) (*sqlcgen.CreateCustomPricingParams, error) {

	var form CustomPricingForm

	if err := c.ShouldBind(&form); err != nil {
		return nil, err
	}

	if form.CostMode == "" {
		form.CostMode = "fixed"
	}

	input := sqlcgen.CreateCustomPricingParams{
		TierBasePricingID: form.TierBasePricingID,
		SubscriptionID:    form.SubscriptionID,
		CustomRateLimit:   form.CustomRateLimit,
		CustomCostPerCall: form.CustomCostPerCall,
		CostMode:          form.CostMode,
	}

	return &input, nil
}

func (h *PricingService) CreateCustomPricingJSONValidation(c *gin.Context) ([]sqlcgen.CreateCustomPricingsParams, error) {

	var inputs []CreateCustomPricingParams

	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var outputs []sqlcgen.CreateCustomPricingsParams

	for _, input := range inputs {
		batchInput := sqlcgen.CreateCustomPricingsParams{
			CustomCostPerCall: input.CustomCostPerCall,
			CustomRateLimit:   int32(input.CustomRateLimit),
			SubscriptionID:    int32(input.SubscriptionID),
			TierBasePricingID: int32(input.TierBasePricingID),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}
