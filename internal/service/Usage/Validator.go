package usage

import (
	"fmt"

	"github.com/bignyap/go-admin/database/sqlcgen"
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
