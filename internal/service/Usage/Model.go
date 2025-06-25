package usage

import (
	"time"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
)

type CreateApiUsageSummaryParams struct {
	StartDate      time.Time `json:"start_date" form:"start_date"`
	EndDate        time.Time `json:"end_date" form:"end_date"`
	TotalCalls     int       `json:"total_calls" form:"total_calls"`
	TotalCost      float64   `json:"total_cost" form:"total_cost"`
	SubscriptionId int       `json:"subscription_id" form:"subscription_id"`
	ApiEndpointId  int       `json:"api_endpoint_id" form:"api_endpoint_id"`
	OrganizationId int       `json:"organization_id" form:"organization_id"`
}

type CreateApiUsageSummaryOutput struct {
	ID int `json:"id"`
	CreateApiUsageSummaryParams
}

type CreateApiUsageSummaryInput interface {
	ToCreateApiUsageSummaryParams() CreateApiUsageSummaryParams
}

type LocalApiUsageSummary struct {
	sqlcgen.ApiUsageSummary
}

type LocalCreateApiUsageSummaryParams struct {
	sqlcgen.CreateApiUsageSummaryParams
}

func (apiSummary LocalApiUsageSummary) ToCreateApiUsageSummaryParams() CreateApiUsageSummaryParams {
	return CreateApiUsageSummaryParams{
		StartDate:      converter.FromUnixTime32(apiSummary.UsageStartDate),
		EndDate:        converter.FromUnixTime32(apiSummary.UsageEndDate),
		TotalCalls:     int(apiSummary.TotalCalls),
		TotalCost:      apiSummary.TotalCost,
		SubscriptionId: int(apiSummary.SubscriptionID),
		ApiEndpointId:  int(apiSummary.ApiEndpointID),
		OrganizationId: int(apiSummary.OrganizationID),
	}
}

func (apiSummary LocalCreateApiUsageSummaryParams) ToCreateApiUsageSummaryParams() CreateApiUsageSummaryParams {
	return CreateApiUsageSummaryParams{
		StartDate:      converter.FromUnixTime32(apiSummary.UsageStartDate),
		EndDate:        converter.FromUnixTime32(apiSummary.UsageEndDate),
		TotalCalls:     int(apiSummary.TotalCalls),
		TotalCost:      apiSummary.TotalCost,
		SubscriptionId: int(apiSummary.SubscriptionID),
		ApiEndpointId:  int(apiSummary.ApiEndpointID),
		OrganizationId: int(apiSummary.OrganizationID),
	}
}

func ToCreateApiUsageSummaryOutput(input sqlcgen.ApiUsageSummary) CreateApiUsageSummaryOutput {
	return CreateApiUsageSummaryOutput{
		ID:                          int(input.OrganizationID),
		CreateApiUsageSummaryParams: LocalApiUsageSummary{input}.ToCreateApiUsageSummaryParams(),
	}
}
