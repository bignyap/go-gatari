package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bignyap/go-admin/internal/common"
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
	if cfg.MemcacheCfg == nil {
		cfg.MemcacheCfg = &memcache.Config{
			DefaultTTL:      cfg.LocalTTL,
			CleanupInterval: 5 * time.Minute,
		}
	}
	memClient := memcache.New(*cfg.MemcacheCfg)

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

func (cc *CacheController) Redis() redis.UniversalClient {
	return cc.redis
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

func (cc *CacheController) SetWithTTL(ctx context.Context, key common.RedisPrefix, val interface{}) error {
	ttl := common.TTLFor(key)

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return cc.redis.Set(ctx, string(key), b, ttl).Err()
}

func (cc *CacheController) Invalidate(ctx context.Context, key string) {
	cc.local.Delete(key)
	if cc.redis != nil {
		cc.redis.Del(ctx, key)
	}
}

// Output looks like:
//
//	map[string]map[string]float64{
//	  "<orgId>:<subId>:<endpointId>": {
//	    "Suffix1": 10.0,
//	    "Suffix2": 0.0,
//	    "Suffix3": 5.0,
//	  },
//	  "<orgId>:<subId>:<endpointId>": {
//	    "Suffix1": 10.0,
//	    "Suffix2": 0.0,
//	    "Suffix3": 5.0,
//	  },
//	}
//
// If the input suffixes are ["Suffix1", "Suffix2", "Suffix3"],
// the output will always contain these keys for each ID, even if no data was found in Redis.
func (cc *CacheController) GetRedisGroupedSnapshot(
	ctx context.Context,
	prefix string,
	suffixes []string,
) map[string]map[string]float64 {

	result := make(map[string]map[string]float64)

	if cc.redis == nil {
		return result
	}

	suffixSet := make(map[string]struct{})
	for _, s := range suffixes {
		suffixSet[s] = struct{}{}
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
			if err != nil {
				continue
			}

			// key = prefix:<logical_id>:<suffix>
			keySuffix := strings.TrimPrefix(key, prefix+":")
			parts := strings.Split(keySuffix, ":")
			if len(parts) < 2 {
				continue
			}

			suffix := parts[len(parts)-1]
			if _, ok := suffixSet[suffix]; !ok {
				continue
			}

			id := strings.Join(parts[:len(parts)-1], ":")

			if _, ok := result[id]; !ok {
				// Initialize all suffixes with 0
				result[id] = make(map[string]float64)
				for _, sfx := range suffixes {
					result[id][sfx] = 0.0
				}
			}

			// Parse float and assign
			var val float64
			fmt.Sscanf(valStr, "%f", &val)
			result[id][suffix] = val
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	// Ensure all suffixes exist per ID (even if no key was found)
	for id := range result {
		for _, sfx := range suffixes {
			if _, ok := result[id][sfx]; !ok {
				result[id][sfx] = 0.0
			}
		}
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
		cc.redis.Del(ctx, prefix+":"+key)
	}
}

func (cc *CacheController) Close() error {
	if closer, ok := cc.redis.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

type PrimitiveWrapper[T any] struct {
	Value T `json:"value"`
}
