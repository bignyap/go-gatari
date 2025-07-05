package cachemanagement

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bignyap/go-admin/internal/common"
)

func StartPeriodicFlush(cm *CacheManagementService, interval time.Duration, stopCh <-chan struct{}) {

	cm.Logger.WithComponent("StartPeriodicFlush").Info("Started")

	go func() {
		for {
			now := time.Now()
			next := now.Truncate(interval).Add(interval)
			wait := time.Until(next)

			timer := time.NewTimer(wait)

			select {
			case <-timer.C:
				ctx := context.Background()
				cm.SyncAggregatedToDB(ctx, string(common.UsagePrefix), func(key string, val map[string]float64) error {
					return cm.IncrementUsageFromCacheKey(ctx, key, val)
				})

			case <-stopCh:
				cm.Logger.WithComponent("StartPeriodicFlush").Info("Stopped")
				return
			}
		}
	}()
}

func HandleShutdown(cm *CacheManagementService, stopFlush chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	close(stopFlush)

	ctx := context.Background()
	cm.SyncAggregatedToDB(ctx, string(common.UsagePrefix), func(key string, val map[string]float64) error {
		return cm.IncrementUsageFromCacheKey(ctx, key, val)
	})

	cm.Logger.WithComponent("HandleShutdown").Info("Gatekeeper shutdown complete")

	os.Exit(0)
}
