package organization

import (
	"time"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
)

type CreateOrgTypeInput struct {
	Name string `json:"name" form:"name" validate:"required,min=1"`
}

type CreateOrgTypeOutput struct {
	ID int `json:"id" form:"id"`
	CreateOrgTypeInput
}

type CreateOrgTypeParams struct {
	Names []string `json:"name" form:"name" validate:"required,dive,required,min=1"`
}

type CreateOrgPermissionParams struct {
	ResourceTypeID int    `json:"resource_type_id" form:"resource_type_id" validate:"required"`
	OrganizationID int    `json:"organization_id" form:"organization_id" validate:"required"`
	PermissionCode string `json:"permission_code" form:"permission_code" validate:"required"`
}

type CreateOrgPermissionOutput struct {
	ID int `json:"id" form:"id"`
	CreateOrgPermissionParams
}

type CreateOrganizationParams struct {
	Name         string    `json:"name" form:"name" validate:"required"`
	CreatedAt    time.Time `json:"created_at" form:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" form:"updated_at"`
	Realm        string    `json:"realm" form:"realm" validate:"required"`
	Country      *string   `json:"country" form:"country"`
	SupportEmail string    `json:"support_email" form:"support_email" validate:"required,email"`
	Active       *bool     `json:"active" form:"active"`
	ReportQ      *bool     `json:"report_q" form:"report_q"`
	Config       *string   `json:"config" form:"config"`
	TypeID       int       `json:"type_id" form:"type_id" validate:"required,min=1"`
}

type CreateOrganizationOutput struct {
	ID int `json:"id" form:"id"`
	CreateOrganizationParams
}

type ListOrganizationOutput struct {
	ID                   int    `json:"id" form:"id"`
	OrganizationTypeName string `json:"type" form:"type"`
	CreateOrganizationParams
}

type ListOrganizationOutputWithCount struct {
	TotalItems int                      `json:"total_items" form:"total_items"`
	Data       []ListOrganizationOutput `json:"data" form:"data"`
}

func ToListOrganizationOutputWithCount(inputs []sqlcgen.ListOrganizationRow) ListOrganizationOutputWithCount {
	var data []ListOrganizationOutput
	for _, input := range inputs {
		data = append(data, ToListOrganizationOutput(input))
	}

	totalItems := 0
	if len(inputs) > 0 {
		totalItems = int(inputs[0].TotalItems)
	}

	return ListOrganizationOutputWithCount{
		Data:       data,
		TotalItems: totalItems,
	}
}

func ToListOrganizationOutput(input sqlcgen.ListOrganizationRow) ListOrganizationOutput {
	return ListOrganizationOutput{
		ID:                   int(input.OrganizationID),
		OrganizationTypeName: input.OrganizationTypeName,
		CreateOrganizationParams: CreateOrganizationParams{
			Name:         input.OrganizationName,
			SupportEmail: input.OrganizationSupportEmail,
			CreatedAt:    converter.FromUnixTime32(input.OrganizationCreatedAt),
			UpdatedAt:    converter.FromUnixTime32(input.OrganizationUpdatedAt),
			Realm:        input.OrganizationRealm,
			Active:       &input.OrganizationActive.Bool,
			ReportQ:      &input.OrganizationReportQ.Bool,
			TypeID:       int(input.OrganizationTypeID),
			Config:       &input.OrganizationConfig.String,
			Country:      &input.OrganizationCountry.String,
		},
	}
}

type UpdateOrganizationParams struct {
	Name           string  `json:"name" form:"name" validate:"required"`
	Realm          string  `json:"realm" form:"realm" validate:"required"`
	Country        *string `json:"country" form:"country"`
	SupportEmail   string  `json:"support_email" form:"support_email" validate:"required,email"`
	Active         *bool   `json:"active" form:"active"`
	ReportQ        *bool   `json:"report_q" form:"report_q"`
	Config         *string `json:"config" form:"config"`
	TypeID         int     `json:"type_id" form:"type_id" validate:"required,min=1"`
	OrganizationID int     `json:"organization_id" form:"organization_id" validate:"required,min=1"`
}
