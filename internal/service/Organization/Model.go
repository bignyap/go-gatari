package organization

import (
	"time"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/utils/misc"
)

type CreateOrgTypeInput struct {
	Name string `json:"name" validate:"required,min=1"`
}

type CreateOrgTypeOutput struct {
	ID int `json:"id"`
	CreateOrgTypeInput
}

type CreateOrgTypeParams struct {
	Names []string `json:"name" validate:"required,dive,required,min=1"`
}

type CreateOrgPermissionParams struct {
	ResourceTypeID int    `json:"resource_type_id" validate:"required"`
	OrganizationID int    `json:"organization_id" validate:"required"`
	PermissionCode string `json:"permission_code" validate:"required"`
}

type CreateOrgPermissionOutput struct {
	ID int `json:"id"`
	CreateOrgPermissionParams
}

type CreateOrganizationParams struct {
	Name         string    `json:"name" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Realm        string    `json:"realm" validate:"required"`
	Country      *string   `json:"country"`
	SupportEmail string    `json:"support_email" validate:"required,email"`
	Active       *bool     `json:"active"`
	ReportQ      *bool     `json:"report_q"`
	Config       *string   `json:"config"`
	TypeID       int       `json:"type_id" validate:"required,min=1"`
}

type CreateOrganizationOutput struct {
	ID int `json:"id"`
	CreateOrganizationParams
}

type ListOrganizationOutput struct {
	ID                   int    `json:"id"`
	OrganizationTypeName string `json:"type"`
	CreateOrganizationParams
}

type ListOrganizationOutputWithCount struct {
	TotalItems int                      `json:"total_items"`
	Data       []ListOrganizationOutput `json:"data"`
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
			CreatedAt:    misc.FromUnixTime32(input.OrganizationCreatedAt),
			UpdatedAt:    misc.FromUnixTime32(input.OrganizationUpdatedAt),
			Realm:        input.OrganizationRealm,
			Active:       &input.OrganizationActive.Bool,
			ReportQ:      &input.OrganizationReportQ.Bool,
			TypeID:       int(input.OrganizationTypeID),
			Config:       &input.OrganizationConfig.String,
			Country:      &input.OrganizationCountry.String,
		},
	}
}
