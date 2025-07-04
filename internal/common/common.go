package common

import (
	"fmt"
	"strconv"
	"strings"
)

type PubSubChannel string

const (
	EndpointCreated      PubSubChannel = "endpoint:created"
	EndpointDeleted      PubSubChannel = "endpoint:deleted"
	OrganizationModified PubSubChannel = "organization:modified"
	SubscriptionModified PubSubChannel = "subscription:modified"
	PricingModified      PubSubChannel = "pricing:modified"
)

type RedisPrefix string

const (
	Usageprefix RedisPrefix = "usage"
)

type RedisKeyPrefix string

const (
	Countprefix           RedisKeyPrefix = "Count"
	Costprefix            RedisKeyPrefix = "Cost"
	TotalCostKeyPrefix    RedisKeyPrefix = "TotalCost"
	TotalCountKeyPrefix   RedisKeyPrefix = "TotalCount"
	UsageKeyPrefix        RedisKeyPrefix = "Usage"
	PricingKeyPrefix      RedisKeyPrefix = "Pricing"
	OrganizationKeyPrefix RedisKeyPrefix = "Organization"
	EndpointKeyPrefix     RedisKeyPrefix = "Endpoint"
	SubscriptionKeyPrefix RedisKeyPrefix = "Subscription"
)

func ParseUsageKey(key string) (orgID, subID, endpointID int32, err error) {

	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return 0, 0, 0, fmt.Errorf("invalid key format")
	}
	org64, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid orgID")
	}
	sub64, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid subscriptionID")
	}
	endpoint64, err := strconv.ParseInt(parts[3], 10, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid endpointID")
	}
	return int32(org64), int32(sub64), int32(endpoint64), nil
}

func RedisKeyFormatter[T ~string | ~int | ~int32 | ~int64](args ...T) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = fmt.Sprint(a)
	}
	return strings.Join(parts, ":")
}

func UsageKey(orgID, subID, endpointID int32) string {
	return RedisKeyFormatter(orgID, subID, endpointID)
}

type EndpointCreatedEvent struct {
	Code   string
	Path   string
	Method string
}

type EndpointDeletedEvent struct {
	Code string
}

func FetchAll[T any](fetchFunc func(offset, batchsize int32) ([]T, error), batchsize int32) ([]T, error) {

	var results []T
	offset := int32(0)

	for {
		items, err := fetchFunc(offset, batchsize)
		if err != nil {
			return nil, err
		}

		results = append(results, items...)

		if int32(len(items)) < batchsize {
			break
		}

		offset += batchsize
	}

	return results, nil
}

func SafeGet[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultVal
}
