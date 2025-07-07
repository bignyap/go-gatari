package usage

import (
	"context"
	"fmt"
	"time"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

type UsageSummaryFilters struct {
	OrgID      *int  `form:"organization_id"`
	SubID      *int  `form:"subscription_id"`
	EndpointID *int  `form:"endpoint_id"`
	StartDate  *int  `form:"start_date"`
	EndDate    *int  `form:"end_date"`
	Limit      int32 `form:"limit" binding:"required" default:"100"`
	Offset     int32 `form:"offset" default:"1"`
	GroupBy    bool  `form:"group_by" default:"false"`
}

func (s *UsageSummaryService) UsageSummaryQueryValidation(c *gin.Context) (UsageSummaryFilters, error) {
	var filters UsageSummaryFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		return UsageSummaryFilters{}, err
	}

	// Validate date range
	if filters.StartDate != nil && filters.EndDate != nil && *filters.StartDate > *filters.EndDate {
		return UsageSummaryFilters{}, fmt.Errorf("start_date cannot be greater than end_date")
	}

	return filters, nil
}

func GetGroupedUsageSummary[T any, P any](
	ctx context.Context,
	queryFunc func(context.Context, P) ([]T, error),
	paramBuilder func(UsageSummaryFilters) P,
	filters UsageSummaryFilters,
) ([]T, error) {
	params := paramBuilder(filters)
	return queryFunc(ctx, params)
}

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
