package caching

import (
	"context"
	"encoding/json"
	"fmt"
)

func GetFromCache[T any](
	ctx context.Context,
	cache *CacheController,
	key string,
	fallback func() (T, error),
) (T, error) {

	var zero T

	val, err := cache.Get(ctx, key, func() (interface{}, error) {
		v, err := fallback()
		if err != nil {
			return nil, err
		}

		fmt.Println("value: ", v)

		// Wrap primitives before serializing
		if isPrimitiveType[T]() {
			return PrimitiveWrapper[T]{Value: v}, nil
		}

		return v, nil
	})
	if err != nil {
		return zero, err
	}

	// Handle primitives
	if isPrimitiveType[T]() {
		// Marshal/unmarshal back to wrapper
		b, err := json.Marshal(val)
		if err != nil {
			return zero, err
		}
		var wrapper PrimitiveWrapper[T]
		if err := json.Unmarshal(b, &wrapper); err != nil {
			return zero, err
		}
		return wrapper.Value, nil
	}

	// Handle structs
	if result, ok := val.(T); ok {
		return result, nil
	}

	// Try converting via JSON
	b, err := json.Marshal(val)
	if err != nil {
		return zero, err
	}
	var typed T
	if err := json.Unmarshal(b, &typed); err != nil {
		return zero, err
	}
	return typed, nil
}

func isPrimitiveType[T any]() bool {
	var t T
	switch any(t).(type) {
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}
