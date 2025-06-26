package caching

import (
	"context"
	"time"

	"github.com/bignyap/go-utilities/memcache"
	"github.com/bignyap/go-utilities/redisclient"
)

type CacheController struct {
	LocalTTL     time.Duration
	RedisTTL     time.Duration
	Serializer   func(interface{}) (string, error)
	Deserializer func(string) (interface{}, error)
}

func NewCacheController(localTTL, redisTTL time.Duration, ser func(interface{}) (string, error), deser func(string) (interface{}, error)) *CacheController {
	return &CacheController{
		LocalTTL:     localTTL,
		RedisTTL:     redisTTL,
		Serializer:   ser,
		Deserializer: deser,
	}
}

func (cc *CacheController) Get(ctx context.Context, key string, fetchFromSource func() (interface{}, error)) (interface{}, error) {
	// Check in-memory cache
	if val, ok := memcache.Get(key); ok {
		return val, nil
	}

	// Check Redis
	rdb, err := redisclient.GetRedisClient()
	if err == nil {
		if strVal, err := rdb.Get(ctx, key).Result(); err == nil {
			val, err := cc.Deserializer(strVal)
			if err == nil {
				memcache.Set(key, val, cc.LocalTTL)
				return val, nil
			}
		}
	}

	// Fetch from source
	val, err := fetchFromSource()
	if err != nil {
		return nil, err
	}

	// Serialize and store in Redis and memcache
	if rdb != nil {
		if strVal, err := cc.Serializer(val); err == nil {
			rdb.Set(ctx, key, strVal, cc.RedisTTL)
		}
	}
	memcache.Set(key, val, cc.LocalTTL)
	return val, nil
}

func (cc *CacheController) Invalidate(key string) {
	memcache.Invalidate(key)
	if rdb, err := redisclient.GetRedisClient(); err == nil {
		rdb.Del(context.Background(), key)
	}
}
