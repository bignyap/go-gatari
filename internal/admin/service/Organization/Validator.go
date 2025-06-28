package organization

import (
	"fmt"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/gin-gonic/gin"
)

func (h *OrganizationService) CreateOrgTypeFormValidator(c *gin.Context) (string, error) {

	var input CreateOrgTypeInput

	if err := c.ShouldBind(&input); err != nil {
		return "", err
	}

	if err := h.Validator.Struct(input); err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}

	return input.Name, nil
}

func (h *OrganizationService) CreateOrgTypeJSONValidation(c *gin.Context) (CreateOrgTypeParams, error) {

	var inputs []CreateOrgTypeInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return CreateOrgTypeParams{}, fmt.Errorf("invalid JSON: %w", err)
	}

	for _, input := range inputs {
		if err := h.Validator.Struct(input); err != nil {
			return CreateOrgTypeParams{}, fmt.Errorf("validation failed: %w", err)
		}
	}

	var names []string
	for _, val := range inputs {
		names = append(names, val.Name)
	}

	return CreateOrgTypeParams{Names: names}, nil
}

func (h *OrganizationService) CreateOrgPermissionJSONValidation(c *gin.Context) ([]sqlcgen.CreateOrgPermissionsParams, error) {

	var inputs []CreateOrgPermissionParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var outputs []sqlcgen.CreateOrgPermissionsParams
	for _, input := range inputs {
		if err := h.Validator.Struct(input); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		outputs = append(outputs, sqlcgen.CreateOrgPermissionsParams{
			OrganizationID: int32(input.OrganizationID),
			ResourceTypeID: int32(input.ResourceTypeID),
			PermissionCode: input.PermissionCode,
		})
	}

	return outputs, nil
}

type CreateOrgPermissionInput struct {
	PermissionCode string `form:"permission_code" binding:"required"`
	ResourceTypeID int32  `form:"resource_type_id"`
	OrganizationID int32  `form:"organization_id"`
}

func (h *OrganizationService) CreateOrgPermissionFormValidator(c *gin.Context) (*sqlcgen.CreateOrgPermissionParams, error) {

	var input CreateOrgPermissionInput

	if err := c.ShouldBind(&input); err != nil {
		return nil, fmt.Errorf("binding failed: %w", err)
	}

	if err := h.Validator.Struct(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &sqlcgen.CreateOrgPermissionParams{
		OrganizationID: input.OrganizationID,
		ResourceTypeID: input.ResourceTypeID,
		PermissionCode: input.PermissionCode,
	}, nil
}

func (h *OrganizationService) ValidateOrgInput(c *gin.Context) (*CreateOrganizationParams, error) {

	var inputs CreateOrganizationParams
	if err := c.ShouldBind(&inputs); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	if err := h.Validator.Struct(inputs); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	return &inputs, nil
}

func (h *OrganizationService) ValidateOrgBatchInput(c *gin.Context) ([]CreateOrganizationParams, error) {
	var inputs []CreateOrganizationParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	for i, input := range inputs {
		if err := h.Validator.Struct(input); err != nil {
			return nil, fmt.Errorf("validation error at index %d: %w", i, err)
		}
	}
	return inputs, nil
}
