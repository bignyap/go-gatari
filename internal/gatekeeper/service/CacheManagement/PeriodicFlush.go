package cachemanagement

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartPeriodicFlush(cm *CacheManagementService, interval time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx := context.Background()

				cm.SyncIncrementalToRedis(ctx, "usage")

				cm.SyncAggregatedToDB(ctx, "usage", func(key string, count int) error {
					return cm.IncrementUsageFromCacheKey(ctx, key, count)
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

	cm.SyncIncrementalToRedis(ctx, "usage")

	cm.SyncAggregatedToDB(ctx, "usage", func(key string, count int) error {
		return cm.IncrementUsageFromCacheKey(ctx, key, count)
	})

	log.Println("Gatekeeper shutdown complete")
	os.Exit(0)
}
