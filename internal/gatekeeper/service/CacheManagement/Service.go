package cachemanagement

import (
	"context"

	"github.com/bignyap/go-utilities/counter"

	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CacheManagementService struct {
	FlushInterval int64

	DB        *sqlcgen.Queries
	Conn      *pgxpool.Pool
	Logger    api.Logger
	Validator *validator.Validate
	Cache     *caching.CacheController

	// New additions
	CounterWorker     *counter.CounterWorker
	RedisSnapshotFunc func(ctx context.Context, prefix string, suffix []string) map[string]map[string]float64
	RedisResetFunc    func(ctx context.Context, prefix string)
}

// Helper to safely inject Redis functions
func NewCacheManagementService(
	flushInterval int64,
	db *sqlcgen.Queries,
	conn *pgxpool.Pool,
	logger api.Logger,
	validator *validator.Validate,
	cache *caching.CacheController,
	redis redis.UniversalClient,
	counterWorker *counter.CounterWorker,
) *CacheManagementService {
	return &CacheManagementService{
		FlushInterval:     flushInterval,
		DB:                db,
		Conn:              conn,
		Logger:            logger,
		Validator:         validator,
		Cache:             cache,
		CounterWorker:     counterWorker,
		RedisSnapshotFunc: cache.GetRedisGroupedSnapshot,
		RedisResetFunc:    cache.ResetRedisValues,
	}
}
