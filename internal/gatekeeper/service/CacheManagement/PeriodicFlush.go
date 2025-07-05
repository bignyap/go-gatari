package cachemanagement

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bignyap/go-admin/internal/common"
)

func StartPeriodicFlush(cm *CacheManagementService, interval time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx := context.Background()

				// Only this node flushes Redis -> DB
				cm.SyncAggregatedToDB(ctx, string(common.UsagePrefix), func(key string, val map[string]float64) error {
					return cm.IncrementUsageFromCacheKey(ctx, key, val)
				})

			case <-stopCh:
				log.Println("Stopping cache flush ticker")
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

	log.Println("Gatekeeper shutdown complete")
	os.Exit(0)
}
