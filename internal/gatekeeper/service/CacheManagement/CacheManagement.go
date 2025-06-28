package cachemanagement

import (
	"context"
	"fmt"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
)

// Push local incremental values to Redis
func (srvc *CacheManagementService) SyncIncrementalToRedis(ctx context.Context, prefix string) {
	data := srvc.Cache.GetAllLocalValues(prefix) // map[string]int
	for key, val := range data {
		err := srvc.Cache.IncrementRedisValue(ctx, key, val)
		if err != nil {
			srvc.Logger.Error(
				fmt.Sprintf("failed to increment redis for key %s", key), err,
			)
		}
	}
	srvc.Cache.ResetLocalValues(prefix)
}

// Push aggregated Redis values to DB
func (srvc *CacheManagementService) SyncAggregatedToDB(ctx context.Context, prefix string, handler func(key string, count int) error) {
	data := srvc.Cache.GetRedisSnapshot(ctx, prefix) // map[string]int
	for key, val := range data {
		if err := handler(key, val); err != nil {
			srvc.Logger.Error(
				fmt.Sprintf("failed handling redis data for key %s", key), err,
			)
		}
	}
	srvc.Cache.ResetRedisValues(ctx, prefix)
}

func (srvc *CacheManagementService) InvalidateLocal(prefix, key string) {
	srvc.Logger.Info(fmt.Sprintf("Invalidating local cache for key: %s", key))
	srvc.Cache.DeleteLocalValue(prefix, key)
}

func (srvc *CacheManagementService) InvalidateRedis(ctx context.Context, prefix, key string) {
	srvc.Logger.Info(fmt.Sprintf("Invalidating redis cache for key: %s", key))
	srvc.Cache.DeleteRedisValue(ctx, prefix, key)
}

func (srvc *CacheManagementService) IncrementUsageFromCacheKey(ctx context.Context, key string, count int) error {

	orgID, endpointID, err := ParseUsageKey(key)
	if err != nil {
		return err
	}

	sub, err := srvc.DB.GetActiveSubscription(ctx, sqlcgen.GetActiveSubscriptionParams{
		OrganizationID: orgID,
		ApiEndpointID:  endpointID,
	})
	if err != nil || !sub.Active.Bool {
		return fmt.Errorf("no active subscription for org=%d endpoint=%d", orgID, endpointID)
	}

	return srvc.DB.IncrementUsage(ctx, sqlcgen.IncrementUsageParams{
		OrganizationID: orgID,
		SubscriptionID: sub.ID,
		ApiEndpointID:  endpointID,
	})
}
