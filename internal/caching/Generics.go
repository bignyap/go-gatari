package caching

import (
	"context"
	"encoding/json"
	"fmt"
)

func GetFromCache[T any](ctx context.Context, cache *CacheController, key string, fallback func() (T, error)) (T, error) {
	var zero T

	// First try cache
	val, err := cache.Get(ctx, key, func() (interface{}, error) {
		v, err := fallback()
		if err != nil {
			return nil, err
		}
		return v, nil
	})
	if err != nil {
		return zero, err
	}

	// Now handle conversion
	switch v := val.(type) {
	case T:
		return v, nil
	case map[string]interface{}:
		// Marshal back to JSON, then unmarshal into T
		b, err := json.Marshal(v)
		if err != nil {
			return zero, fmt.Errorf("marshal map failed: %w", err)
		}
		var result T
		if err := json.Unmarshal(b, &result); err != nil {
			return zero, fmt.Errorf("unmarshal to %T failed: %w", result, err)
		}
		return result, nil
	default:
		return zero, fmt.Errorf("type assertion to %T failed, got %T", zero, val)
	}
}
