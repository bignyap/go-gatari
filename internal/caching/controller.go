package caching

import (
	"context"
	"time"

	"github.com/bignyap/go-utilities/memcache"
	"github.com/bignyap/go-utilities/redisclient"
	"github.com/redis/go-redis/v9"
)

type CacheController struct {
	local       *memcache.Client
	redis       redis.UniversalClient
	localTTL    time.Duration
	redisTTL    time.Duration
	serialize   func(interface{}) (string, error)
	deserialize func(string) (interface{}, error)
}

type CacheControllerConfig struct {
	LocalTTL     time.Duration
	RedisTTL     time.Duration
	MemcacheCfg  *memcache.Config
	RedisCfg     *redisclient.RedisConfig
	Serializer   func(interface{}) (string, error)
	Deserializer func(string) (interface{}, error)
}

func NewCacheController(ctx context.Context, cfg CacheControllerConfig) (*CacheController, error) {
	memCfg := cfg.MemcacheCfg
	if memCfg == nil {
		memCfg = &memcache.Config{
			DefaultTTL:      cfg.LocalTTL,
			CleanupInterval: 5 * time.Minute,
		}
	}
	memClient := memcache.New(*memCfg)

	var redisClient redis.UniversalClient
	var err error
	if cfg.RedisCfg != nil {
		redisClient, err = redisclient.New(ctx, *cfg.RedisCfg)
		if err != nil {
			return nil, err
		}
	}

	return &CacheController{
		local:       memClient,
		redis:       redisClient,
		localTTL:    cfg.LocalTTL,
		redisTTL:    cfg.RedisTTL,
		serialize:   cfg.Serializer,
		deserialize: cfg.Deserializer,
	}, nil
}

func (cc *CacheController) Get(ctx context.Context, key string, fetch func() (interface{}, error)) (interface{}, error) {
	// In-memory
	if val, ok := cc.local.Get(key); ok {
		return val, nil
	}

	// Redis
	if cc.redis != nil {
		if strVal, err := cc.redis.Get(ctx, key).Result(); err == nil {
			val, err := cc.deserialize(strVal)
			if err == nil {
				cc.local.Set(key, val, cc.localTTL)
				return val, nil
			}
		}
	}

	// Fallback
	val, err := fetch()
	if err != nil {
		return nil, err
	}

	// Save to both
	cc.local.Set(key, val, cc.localTTL)

	if cc.redis != nil {
		if strVal, err := cc.serialize(val); err == nil {
			cc.redis.Set(ctx, key, strVal, cc.redisTTL)
		}
	}

	return val, nil
}

func (cc *CacheController) Set(ctx context.Context, key string, val interface{}) error {
	cc.local.Set(key, val, cc.localTTL)

	if cc.redis != nil {
		strVal, err := cc.serialize(val)
		if err != nil {
			return err
		}
		return cc.redis.Set(ctx, key, strVal, cc.redisTTL).Err()
	}

	return nil
}

func (cc *CacheController) Invalidate(ctx context.Context, key string) {
	cc.local.Delete(key)
	if cc.redis != nil {
		cc.redis.Del(ctx, key)
	}
}
