package usage

import (
	"context"
	"fmt"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func CreateApiUsageSummaryJSONValidation(c *gin.Context) ([]sqlcgen.CreateApiUsageSummariesParams, error) {

	var inputs []CreateApiUsageSummaryParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var outputs []sqlcgen.CreateApiUsageSummariesParams

	for _, input := range inputs {
		batchInput := sqlcgen.CreateApiUsageSummariesParams{
			UsageStartDate: int32(*converter.TimePtrToUnixInt(&input.StartDate)),
			UsageEndDate:   int32(*converter.TimePtrToUnixInt(&input.EndDate)),
			TotalCalls:     int32(input.TotalCalls),
			TotalCost:      input.TotalCost,
			SubscriptionID: int32(input.SubscriptionId),
			OrganizationID: int32(input.OrganizationId),
			ApiEndpointID:  int32(input.ApiEndpointId),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

func (s *UsageSummaryService) UsageSummaryQueryValidation(c *gin.Context) (UsageSummaryFilterQueryParams, error) {

	var filters UsageSummaryFilterQueryParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		return UsageSummaryFilterQueryParams{}, err
	}

	if filters.StartDate != nil && filters.EndDate != nil && *filters.StartDate > *filters.EndDate {
		return UsageSummaryFilterQueryParams{}, fmt.Errorf("start_date cannot be greater than end_date")
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
