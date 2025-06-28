package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-utilities/redisclient"
)

func LoadRedisController() (*caching.CacheController, error) {

	serializer := func(val interface{}) (string, error) {
		b, err := json.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("serialize error: %w", err)
		}
		return string(b), nil
	}

	deserializer := func(data string) (interface{}, error) {
		var out interface{}
		if err := json.Unmarshal([]byte(data), &out); err != nil {
			return nil, fmt.Errorf("deserialize error: %w", err)
		}
		return out, nil
	}

	cfg := &redisclient.RedisConfig{
		Addr:     getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       getEnvIntOrDefault("REDIS_DB", 0),
	}

	return caching.NewCacheController(context.Background(), caching.CacheControllerConfig{
		LocalTTL:     5 * time.Minute,
		RedisTTL:     30 * time.Minute,
		Serializer:   serializer,
		Deserializer: deserializer,
		RedisCfg:     cfg,
	})
}

// Optional reusable helpers
func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvIntOrDefault(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i
		}
	}
	return fallback
}
