package subscription

import (
	"fmt"
	"time"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *SubscriptionService) CreateSubscriptionInBatchValidation(c *gin.Context) ([]CreateSubscriptionParams, error) {

	var inputs []CreateSubscriptionParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return inputs, nil
}

func (h *SubscriptionService) CreateSubscriptionValidation(c *gin.Context) (*CreateSubscriptionParams, error) {

	var input CreateSubscriptionParams

	if err := c.ShouldBindWith(&input, binding.Form); err != nil {
		return nil, fmt.Errorf("binding error: %w", err)
	}

	// Manual parsing
	if input.StartDateRaw != "" {
		t := &converter.TimeOrDate{}
		if err := t.UnmarshalText([]byte(input.StartDateRaw)); err != nil {
			return nil, fmt.Errorf("invalid start_date: %w", err)
		}
		input.StartDate = t
	}

	if input.ExpiryDateRaw != "" {
		t := &converter.TimeOrDate{}
		if err := t.UnmarshalText([]byte(input.ExpiryDateRaw)); err != nil {
			return nil, fmt.Errorf("invalid expiry_date: %w", err)
		}
		input.ExpiryDate = t
	}

	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if input.StartDate == nil || input.StartDate.IsZero() {
		input.StartDate = &converter.TimeOrDate{Time: now}
	}

	return &input, nil
}

func (h *SubscriptionService) CreateSubscriptionTierValidation(c *gin.Context) (*CreateSubTierParams, error) {

	var input CreateSubTierParams
	if err := c.ShouldBind(&input); err != nil {
		return nil, fmt.Errorf("invalid inputs: %w", err)
	}

	return &input, nil
}

func (h *SubscriptionService) CreateSubscriptionTierJSONValidation(c *gin.Context) ([]sqlcgen.CreateSubscriptionTiersParams, error) {

	var inputs []CreateSubTierParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var outputs []sqlcgen.CreateSubscriptionTiersParams
	currentTime := int32(time.Now().Unix())

	for i, input := range inputs {
		if err := h.Validator.Struct(input); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %s", i, err.Error())
		}

		var description pgtype.Text
		if input.Description != nil {
			description.String = *input.Description
			description.Valid = true
		} else {
			description.Valid = false
		}

		outputs = append(outputs, sqlcgen.CreateSubscriptionTiersParams{
			TierName:        input.Name,
			TierDescription: description,
			TierCreatedAt:   currentTime,
			TierUpdatedAt:   currentTime,
		})
	}

	return outputs, nil
}
