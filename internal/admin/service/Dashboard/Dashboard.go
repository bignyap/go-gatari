package dashboard

import (
	"context"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
)

func (s *DashboardService) GetDashboardCounts(ctx context.Context) ([]sqlcgen.DashboardSummaryView, error) {

	dashboardCounts, err := s.DB.GetDashboardSummary(ctx)
	if err != nil {
		return []sqlcgen.DashboardSummaryView{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the counts",
			err,
		)
	}

	if len(dashboardCounts) == 0 {
		return []sqlcgen.DashboardSummaryView{}, nil
	}

	return dashboardCounts, nil
}

func (s *DashboardService) GetDashboardUsage(
	ctx context.Context,
	input sqlcgen.GetTotalCallsGroupedByOrgAndTimeBucketParams,
) (
	[]sqlcgen.GetTotalCallsGroupedByOrgAndTimeBucketRow, error,
) {

	usages, err := s.DB.GetTotalCallsGroupedByOrgAndTimeBucket(ctx, input)
	if err != nil {
		return []sqlcgen.GetTotalCallsGroupedByOrgAndTimeBucketRow{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the counts",
			err,
		)
	}

	if len(usages) == 0 {
		return []sqlcgen.GetTotalCallsGroupedByOrgAndTimeBucketRow{}, nil
	}

	return usages, nil
}
