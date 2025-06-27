package billing

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
)

func (h *BillingService) CreateBillingHistoryJSONValidation(c *gin.Context) ([]sqlcgen.CreateBillingHistoriesParams, error) {

	var inputs []CreateBillingHistoryParams
	if err := c.ShouldBindJSON(&inputs); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	currentTime := time.Now()
	unixNow := int32(currentTime.Unix())

	var outputs []sqlcgen.CreateBillingHistoriesParams
	for _, input := range inputs {
		input.CreatedAt = currentTime

		if err := h.Validator.Struct(input); err != nil {
			return nil, fmt.Errorf("validation error: %w", err)
		}

		output := sqlcgen.CreateBillingHistoriesParams{
			BillingStartDate: int32(input.StartDate.Unix()),
			BillingEndDate:   int32(input.EndDate.Unix()),
			TotalAmountDue:   input.TotalAmountDue,
			TotalCalls:       int32(input.TotalCalls),
			PaymentStatus:    input.PaymentStatus,
			PaymentDate:      timePtrToInt4(input.PaymentDate),
			BillingCreatedAt: unixNow,
			SubscriptionID:   int32(input.SubscriptionId),
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

func int4ToTimePtr(v pgtype.Int4) *time.Time {
	if v.Valid {
		t := time.Unix(int64(v.Int32), 0)
		return &t
	}
	return nil
}

func timePtrToInt4(t *time.Time) pgtype.Int4 {
	if t == nil {
		return pgtype.Int4{
			Int32: 0,
			Valid: false,
		}
	}
	return pgtype.Int4{
		Int32: int32(t.Unix()),
		Valid: true,
	}
}
