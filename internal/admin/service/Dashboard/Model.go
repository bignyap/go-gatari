package dashboard

import "time"

type DashboardUsageFilterQueryParams struct {
	OrgID      *int       `form:"organization_id"`
	BucketSize int        `form:"bucket_size"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
}
