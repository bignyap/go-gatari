package billing

import (
	"time"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/utils/misc"
)

type CreateBillingHistoryParams struct {
	StartDate      time.Time  `json:"start_date" validate:"required"`
	EndDate        time.Time  `json:"end_date" validate:"required,gtfield=StartDate"`
	TotalAmountDue float64    `json:"total_amount_due" validate:"required,gte=0"`
	TotalCalls     int        `json:"total_calls" validate:"required,gte=0"`
	PaymentStatus  string     `json:"payment_status" validate:"required,oneof=paid pending failed"`
	PaymentDate    *time.Time `json:"payment_date" validate:"omitempty"`
	CreatedAt      time.Time  `json:"created_at" validate:"required"`
	SubscriptionId int        `json:"subscription_id" validate:"required,gt=0"`
}

type CreateBillingHistoryOutput struct {
	ID int `json:"id"`
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
		StartDate:      misc.FromUnixTime32(billingHistory.BillingStartDate),
		EndDate:        misc.FromUnixTime32(billingHistory.BillingEndDate),
		PaymentDate:    int4ToTimePtr(billingHistory.PaymentDate),
		CreatedAt:      misc.FromUnixTime32(billingHistory.BillingCreatedAt),
		TotalCalls:     int(billingHistory.TotalCalls),
		TotalAmountDue: billingHistory.TotalAmountDue,
		PaymentStatus:  billingHistory.PaymentStatus,
		SubscriptionId: int(billingHistory.SubscriptionID),
	}
}

func (billingHistory LocalCreateBillingHistoryParams) ToCreateBillingHistoryParams() CreateBillingHistoryParams {
	return CreateBillingHistoryParams{
		StartDate:      misc.FromUnixTime32(billingHistory.BillingStartDate),
		EndDate:        misc.FromUnixTime32(billingHistory.BillingEndDate),
		PaymentDate:    int4ToTimePtr(billingHistory.PaymentDate),
		CreatedAt:      misc.FromUnixTime32(billingHistory.BillingCreatedAt),
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
