package billing

import (
	"time"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
)

type CreateBillingHistoryParams struct {
	StartDate      time.Time  `json:"start_date" form:"start_date" validate:"required"`
	EndDate        time.Time  `json:"end_date" form:"end_date" validate:"required,gtfield=StartDate"`
	TotalAmountDue float64    `json:"total_amount_due" form:"total_amount_due" validate:"required,gte=0"`
	TotalCalls     int        `json:"total_calls" form:"total_calls" validate:"required,gte=0"`
	PaymentStatus  string     `json:"payment_status" form:"payment_status" validate:"required,oneof=paid pending failed"`
	PaymentDate    *time.Time `json:"payment_date" form:"payment_date" validate:"omitempty"`
	CreatedAt      time.Time  `json:"created_at" form:"created_at" validate:"required"`
	SubscriptionId int        `json:"subscription_id" form:"subscription_id" validate:"required,gt=0"`
}

type CreateBillingHistoryOutput struct {
	ID int `json:"id" form:"id"`
	CreateBillingHistoryParams
}

type LocalCreateBillingHistory struct {
	sqlcgen.BillingHistory
}

type LocalCreateBillingHistoryParams struct {
	sqlcgen.CreateBillingHistoryParams
}

func (billingHistory LocalCreateBillingHistory) ToCreateBillingHistoryParams() CreateBillingHistoryParams {
	return CreateBillingHistoryParams{
		StartDate:      converter.FromUnixTime32(billingHistory.BillingStartDate),
		EndDate:        converter.FromUnixTime32(billingHistory.BillingEndDate),
		PaymentDate:    int4ToTimePtr(billingHistory.PaymentDate),
		CreatedAt:      converter.FromUnixTime32(billingHistory.BillingCreatedAt),
		TotalCalls:     int(billingHistory.TotalCalls),
		TotalAmountDue: billingHistory.TotalAmountDue,
		PaymentStatus:  billingHistory.PaymentStatus,
		SubscriptionId: int(billingHistory.SubscriptionID),
	}
}

func (billingHistory LocalCreateBillingHistoryParams) ToCreateBillingHistoryParams() CreateBillingHistoryParams {
	return CreateBillingHistoryParams{
		StartDate:      converter.FromUnixTime32(billingHistory.BillingStartDate),
		EndDate:        converter.FromUnixTime32(billingHistory.BillingEndDate),
		PaymentDate:    int4ToTimePtr(billingHistory.PaymentDate),
		CreatedAt:      converter.FromUnixTime32(billingHistory.BillingCreatedAt),
		TotalCalls:     int(billingHistory.TotalCalls),
		TotalAmountDue: billingHistory.TotalAmountDue,
		PaymentStatus:  billingHistory.PaymentStatus,
		SubscriptionId: int(billingHistory.SubscriptionID),
	}
}

func ToCreateBillingHistoryOutput(input sqlcgen.BillingHistory) CreateBillingHistoryOutput {
	return CreateBillingHistoryOutput{
		ID:                         int(input.BillingID),
		CreateBillingHistoryParams: LocalCreateBillingHistory{input}.ToCreateBillingHistoryParams(),
	}
}
