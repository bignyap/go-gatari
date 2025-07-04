package cachemanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
)

func (srvc *CacheManagementService) SyncAggregatedToDB(ctx context.Context, prefix string, handler func(key string, val map[string]float64) error) {

	data := srvc.RedisSnapshotFunc(ctx, prefix, []string{"Count", "Cost"})

	for key, val := range data {
		if err := handler(key, val); err != nil {
			srvc.Logger.Error(fmt.Sprintf("failed handling redis data for key %s", key), err)
		}
	}

	if srvc.RedisResetFunc != nil {
		srvc.RedisResetFunc(ctx, prefix)
	}
}

func (srvc *CacheManagementService) IncrementUsageFromCacheKey(ctx context.Context, key string, val map[string]float64) error {

	orgID, subID, endpointID, err := common.ParseUsageKey(key)
	if err != nil {
		return err
	}

	_, err = srvc.DB.CreateApiUsageSummary(ctx, sqlcgen.CreateApiUsageSummaryParams{
		UsageStartDate: int32(time.Now().Unix()),
		UsageEndDate:   int32(time.Now().Unix()),
		TotalCalls:     int32(common.SafeGet(val, "Count", 0)),
		TotalCost:      common.SafeGet(val, "Cost", 0.0),
		SubscriptionID: subID,
		ApiEndpointID:  endpointID,
		OrganizationID: orgID,
	})
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			fmt.Sprintf("failed to increment usage for org=%d sub=%d endpoint=%d", orgID, subID, endpointID),
			err,
		)
	}

	return nil

}
