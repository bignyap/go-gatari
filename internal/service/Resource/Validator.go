package resource

import (
	"fmt"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *ResourceService) ValidateRegisterInput(c *gin.Context) (*RegisterEndpointParams, error) {

	var input RegisterEndpointParams
	if err := c.ShouldBind(&input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := h.Validator.Struct(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &input, nil
}

func (h *ResourceService) ValidateRegisterBatchInput(c *gin.Context) ([]RegisterEndpointParams, error) {

	var inputs []RegisterEndpointParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return []RegisterEndpointParams{}, fmt.Errorf("invalid JSON: %w", err)
	}

	for i, input := range inputs {
		if err := h.Validator.Struct(input); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	return inputs, nil
}

func (h *ResourceService) CreateResourceTypeFormValidator(c *gin.Context) (*sqlcgen.CreateResourceTypeParams, error) {

	var input CreateResourceTypeParams
	if err := c.ShouldBind(&input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Validate struct
	if err := h.Validator.Struct(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert description to pgtype.Text
	desc := pgtype.Text{Valid: false}
	if input.Description != nil {
		desc = pgtype.Text{String: *input.Description, Valid: true}
	}

	// Build and return final params
	return &sqlcgen.CreateResourceTypeParams{
		ResourceTypeName:        input.Name,
		ResourceTypeCode:        input.Code,
		ResourceTypeDescription: desc,
	}, nil
}

func (h *ResourceService) CreateResourceTypeJSONValidation(c *gin.Context) ([]sqlcgen.CreateResourceTypesParams, error) {

	var inputs []CreateResourceTypeParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return []sqlcgen.CreateResourceTypesParams{}, fmt.Errorf("invalid JSON: %w", err)
	}

	var outputs []sqlcgen.CreateResourceTypesParams
	for _, input := range inputs {
		if err := h.Validator.Struct(input); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}

		desc := pgtype.Text{Valid: false}
		if input.Description != nil {
			desc = pgtype.Text{String: *input.Description, Valid: true}
		}

		outputs = append(outputs, sqlcgen.CreateResourceTypesParams{
			ResourceTypeName:        input.Name,
			ResourceTypeCode:        input.Code,
			ResourceTypeDescription: desc,
		})
	}

	return outputs, nil
}
