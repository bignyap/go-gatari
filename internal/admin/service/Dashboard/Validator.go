package dashboard

import (
	"fmt"
	"time"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func (s *DashboardService) DashboardUsageQueryValidation(c *gin.Context) (*sqlcgen.GetTotalCallsGroupedByOrgAndTimeBucketParams, error) {
	var filters DashboardUsageFilterQueryParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		return nil, err
	}

	now := time.Now()
	location := now.Location()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), location)

	// Set defaults
	if filters.StartDate == nil {
		filters.StartDate = &startOfDay
	}
	if filters.EndDate == nil {
		filters.EndDate = &endOfDay
	}

	if filters.StartDate.After(*filters.EndDate) {
		return nil, fmt.Errorf("start_date cannot be greater than end_date")
	}

	// Default bucket size
	if filters.BucketSize == 0 {
		filters.BucketSize = 1
	}

	// Validate
	if err := s.Validator.Struct(filters); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert to *int
	startUnix := int(filters.StartDate.Unix())
	endUnix := int(filters.EndDate.Unix())

	return &sqlcgen.GetTotalCallsGroupedByOrgAndTimeBucketParams{
		BucketSize: int32(filters.BucketSize),
		StartDate:  converter.ToPgInt4(&startUnix),
		EndDate:    converter.ToPgInt4(&endUnix),
		OrgID:      converter.ToPgInt4(filters.OrgID),
	}, nil
}
