package subscription

import (
	"time"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
)

const (
	BillingIntervalMonthly = "monthly"
	BillingIntervalYearly  = "yearly"
	BillingIntervalOnce    = "once"

	BillingModelFlat   = "flat"
	BillingModelUsage  = "usage"
	BillingModelHybrid = "hybrid"

	QuotaResetMonthly = "monthly"
	QuotaResetYearly  = "yearly"
	QuotaResetTotal   = "total"
)

type CreateSubscriptionParams struct {
	Name               string                `json:"name" form:"name" validate:"required"`
	Type               string                `json:"type" form:"type" validate:"required"`
	CreatedAt          time.Time             `json:"created_at" form:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at" form:"updated_at"`
	StartDate          *converter.TimeOrDate `json:"start_date" form:"-"`
	StartDateRaw       string                `form:"start_date" json:"-"`
	APILimit           *int                  `json:"api_limit" form:"api_limit"`
	ExpiryDate         *converter.TimeOrDate `json:"expiry_date" form:"-"`
	ExpiryDateRaw      string                `form:"expiry_date" json:"-"`
	Description        *string               `json:"description" form:"description"`
	Status             *bool                 `json:"status" form:"status"`
	OrganizationID     int                   `json:"organization_id" form:"organization_id" validate:"required"`
	SubscriptionTierID int                   `json:"subscription_tier_id" form:"subscription_tier_id" validate:"required"`
	BillingInterval    *string               `json:"billing_interval" form:"billing_interval" validate:"required,oneof=monthly yearly once"`
	BillingModel       *string               `json:"billing_model" form:"billing_model" validate:"required,oneof=flat usage hybrid"`
	QuotaResetInterval *string               `json:"quota_reset_interval" form:"quota_reset_interval" validate:"required,oneof=monthly yearly total"`
}

type CreateSubscriptionOutput struct {
	ID int `json:"id"`
	CreateSubscriptionParams
}

type UpdateSubscriptionParams struct {
	Name               string                `json:"name" form:"name" validate:"required"`
	StartDate          *converter.TimeOrDate `json:"start_date" form:"-"`
	StartDateRaw       string                `form:"start_date" json:"-"`
	APILimit           *int                  `json:"api_limit" form:"api_limit"`
	ExpiryDate         *converter.TimeOrDate `json:"expiry_date" form:"-"`
	ExpiryDateRaw      string                `form:"expiry_date" json:"-"`
	Description        *string               `json:"description" form:"description"`
	Status             *bool                 `json:"status" form:"status"`
	OrganizationID     int                   `json:"organization_id" form:"organization_id" validate:"required"`
	SubscriptionTierID int                   `json:"subscription_tier_id" form:"subscription_tier_id" validate:"required"`
	BillingInterval    *string               `json:"billing_interval" form:"billing_interval" validate:"required,oneof=monthly yearly once"`
	BillingModel       *string               `json:"billing_model" form:"billing_model" validate:"required,oneof=flat usage hybrid"`
	QuotaResetInterval *string               `json:"quota_reset_interval" form:"quota_reset_interval" validate:"required,oneof=monthly yearly total"`
	SubscriptionID     int                   `json:"subscription_id" form:"subscription_id" validate:"required"`
}

type ListSubscriptionOutput struct {
	ID       int    `json:"id"`
	TierName string `json:"tier_name"`
	CreateSubscriptionParams
}

type ListSubscriptionOutputWithCount struct {
	TotalItems int                      `json:"total_items"`
	Data       []ListSubscriptionOutput `json:"data"`
}

func ToListSubscriptionOutput(input sqlcgen.ListSubscriptionRow) ListSubscriptionOutput {
	startDate := time.Unix(int64(input.SubscriptionStartDate), 0)
	expiryDate := converter.FromPgInt4TimePtr(input.SubscriptionExpiryDate)
	return ListSubscriptionOutput{
		ID:       int(input.SubscriptionID),
		TierName: input.TierName,
		CreateSubscriptionParams: CreateSubscriptionParams{
			Name:               input.SubscriptionName,
			Type:               input.SubscriptionType,
			CreatedAt:          time.Unix(int64(input.SubscriptionCreatedDate), 0),
			UpdatedAt:          time.Unix(int64(input.SubscriptionUpdatedDate), 0),
			StartDate:          converter.ToTimeOrDatePtr(&startDate),
			APILimit:           converter.FromPgInt4Ptr(input.SubscriptionApiLimit),
			ExpiryDate:         converter.ToTimeOrDatePtr(expiryDate),
			Description:        converter.FromPgText(input.SubscriptionDescription),
			Status:             converter.FromPgBool(input.SubscriptionStatus),
			OrganizationID:     int(input.OrganizationID),
			SubscriptionTierID: int(input.SubscriptionTierID),
			BillingInterval:    converter.FromPgText(input.SubscriptionBillingInterval),
			BillingModel:       converter.FromPgText(input.SubscriptionBillingModel),
			QuotaResetInterval: converter.FromPgText(input.SubscriptionQuotaResetInterval),
		},
	}
}

func ToListSubscriptionOutputWithCount(inputs []sqlcgen.ListSubscriptionRow) ListSubscriptionOutputWithCount {
	var data []ListSubscriptionOutput
	for _, input := range inputs {
		data = append(data, ToListSubscriptionOutput(input))
	}

	totalItems := 0
	if len(inputs) > 0 {
		totalItems = int(inputs[0].TotalItems)
	}

	return ListSubscriptionOutputWithCount{
		TotalItems: totalItems,
		Data:       data,
	}
}

type CreateSubTierParams struct {
	Name        string    `json:"name" form:"name" validate:"required"`
	Description *string   `json:"description" form:"description"`
	CreatedAt   time.Time `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at"`
}

type CreateSubTierOutput struct {
	ID       int  `json:"id"`
	Archived bool `json:"archived"`
	CreateSubTierParams
}

type CreateSubTierOutputWithCount struct {
	TotalItems int                   `json:"total_items"`
	Data       []CreateSubTierOutput `json:"data"`
}
