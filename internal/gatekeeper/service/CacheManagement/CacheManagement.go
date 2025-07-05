package cachemanagement

import (
	"context"
	"fmt"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
)

func (srvc *CacheManagementService) SyncAggregatedToDB(ctx context.Context, prefix string, handler func(key string, val map[string]float64) error) {

	data := srvc.RedisSnapshotFunc(ctx, prefix, []string{"count", "cost"})
	var hasError bool

	for key, val := range data {
		if err := handler(key, val); err != nil {
			srvc.Logger.Error(fmt.Sprintf("failed handling redis data for key %s", key), err)
			hasError = true
		}
	}

	if !hasError && srvc.RedisResetFunc != nil {
		srvc.RedisResetFunc(ctx, prefix)
	}
}

func (srvc *CacheManagementService) IncrementUsageFromCacheKey(ctx context.Context, key string, val map[string]float64) error {

	orgID, subID, endpointID, timestamp, err := common.ParseUsageKey(key)
	if err != nil {
		return err
	}

	_, err = srvc.DB.CreateApiUsageSummary(ctx, sqlcgen.CreateApiUsageSummaryParams{
		UsageStartDate: timestamp - int32(srvc.FlushInterval),
		UsageEndDate:   timestamp,
		TotalCalls:     int32(common.SafeGet(val, "count", 0)),
		TotalCost:      common.SafeGet(val, "cost", 0.0),
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
