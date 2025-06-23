package subscription

import (
	"time"

	"github.com/bignyap/go-admin/database/sqlcgen"
)

type CreateSubscriptionParams struct {
	Name               string     `json:"name" validate:"required"`
	Type               string     `json:"type" validate:"required"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	StartDate          time.Time  `json:"start_date"`
	APILimit           *int       `json:"api_limit"`
	ExpiryDate         *time.Time `json:"expiry_date"`
	Description        *string    `json:"description"`
	Status             *bool      `json:"status"`
	OrganizationID     int        `json:"organization_id" validate:"required"`
	SubscriptionTierID int        `json:"subscription_tier_id" validate:"required"`
}

type CreateSubscriptionOutput struct {
	ID int `json:"id"`
	CreateSubscriptionParams
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
	return ListSubscriptionOutput{
		ID:       int(input.SubscriptionID),
		TierName: input.TierName,
		CreateSubscriptionParams: CreateSubscriptionParams{
			Name:               input.SubscriptionName,
			Type:               input.SubscriptionType,
			CreatedAt:          time.Unix(int64(input.SubscriptionCreatedDate), 0),
			UpdatedAt:          time.Unix(int64(input.SubscriptionUpdatedDate), 0),
			StartDate:          time.Unix(int64(input.SubscriptionStartDate), 0),
			APILimit:           fromPgInt4Ptr(input.SubscriptionApiLimit),
			ExpiryDate:         fromPgInt4TimePtr(input.SubscriptionExpiryDate),
			Description:        fromPgText(input.SubscriptionDescription),
			Status:             fromPgBool(input.SubscriptionStatus),
			OrganizationID:     int(input.OrganizationID),
			SubscriptionTierID: int(input.SubscriptionTierID),
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
	Name        string    `json:"name" validate:"required"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
