package pricing

import (
	"fmt"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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
			BaseRateLimit:      toPgInt4(input.BaseRateLimit),
			ApiEndpointID:      int32(input.ApiEndpointId),
			SubscriptionTierID: int32(input.SubscriptionTierID),
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}

func toPgInt4(value *int) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*value), Valid: true}
}

type CustomPricingForm struct {
	TierBasePricingID int32   `form:"tier_base_pricing_id" binding:"required"`
	SubscriptionID    int32   `form:"subscription_id" binding:"required"`
	CustomRateLimit   int32   `form:"custom_rate_limit" binding:"required"`
	CustomCostPerCall float64 `form:"custom_cost_per_call" binding:"required"`
}

func (h *PricingService) CreateCustomPricingFormValidator(c *gin.Context) (*sqlcgen.CreateCustomPricingParams, error) {

	var form CustomPricingForm

	if err := c.ShouldBind(&form); err != nil {
		return nil, err
	}

	input := sqlcgen.CreateCustomPricingParams{
		TierBasePricingID: form.TierBasePricingID,
		SubscriptionID:    form.SubscriptionID,
		CustomRateLimit:   form.CustomRateLimit,
		CustomCostPerCall: form.CustomCostPerCall,
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
