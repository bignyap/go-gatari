package caching

import (
	"context"
	"fmt"
	"strings"
	"sync"
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

	mu     sync.Mutex
	counts map[string]map[string]int // prefix -> key -> count
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
		counts:      make(map[string]map[string]int),
	}, nil
}

func (cc *CacheController) Get(ctx context.Context, key string, fetch func() (interface{}, error)) (interface{}, error) {
	if val, ok := cc.local.Get(key); ok {
		return val, nil
	}
	if cc.redis != nil {
		if strVal, err := cc.redis.Get(ctx, key).Result(); err == nil {
			val, err := cc.deserialize(strVal)
			if err == nil {
				cc.local.Set(key, val, cc.localTTL)
				return val, nil
			}
		}
	}
	val, err := fetch()
	if err != nil {
		return nil, err
	}
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

// ---- Incremental Counter Extensions ----

func (cc *CacheController) IncrementLocalValue(prefix, key string, delta int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if _, ok := cc.counts[prefix]; !ok {
		cc.counts[prefix] = make(map[string]int)
	}
	cc.counts[prefix][key] += delta
}

func (cc *CacheController) GetAllLocalValues(prefix string) map[string]int {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	copy := make(map[string]int)
	for k, v := range cc.counts[prefix] {
		copy[k] = v
	}
	return copy
}

func (cc *CacheController) ResetLocalValues(prefix string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.counts[prefix] = make(map[string]int)
}

func (cc *CacheController) DeleteLocalValue(prefix, key string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if _, ok := cc.counts[prefix]; ok {
		delete(cc.counts[prefix], key)
	}
}

func (cc *CacheController) IncrementRedisValue(ctx context.Context, key string, delta int) error {
	if cc.redis == nil {
		return nil
	}
	return cc.redis.IncrBy(ctx, key, int64(delta)).Err()
}

func (cc *CacheController) GetRedisSnapshot(ctx context.Context, prefix string) map[string]int {
	result := make(map[string]int)
	if cc.redis == nil {
		return result
	}

	var cursor uint64
	pattern := prefix + ":*"

	for {
		keys, nextCursor, err := cc.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		for _, key := range keys {
			valStr, err := cc.redis.Get(ctx, key).Result()
			if err == nil {
				var val int
				fmt.Sscanf(valStr, "%d", &val)
				trimmedKey := strings.TrimPrefix(key, prefix+":")
				result[trimmedKey] = val
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return result
}

func (cc *CacheController) ResetRedisValues(ctx context.Context, prefix string) {
	if cc.redis == nil {
		return
	}

	var cursor uint64
	pattern := prefix + ":*"

	for {
		keys, nextCursor, err := cc.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		if len(keys) > 0 {
			cc.redis.Del(ctx, keys...)
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
}

func (cc *CacheController) DeleteRedisValue(ctx context.Context, prefix, key string) {
	if cc.redis != nil {
		fullKey := prefix + ":" + key
		cc.redis.Del(ctx, fullKey)
	}
}

func (c *CacheController) Close() error {
	if closer, ok := c.redis.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}
