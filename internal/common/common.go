package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

type TimeInterval int64

const (
	DefaultRedisFlushInterval    TimeInterval = 10
	DefaultMemcacheFlushInterval TimeInterval = 2
)

type PubSubChannel string

const (
	EndpointCreated       PubSubChannel = "endpoint:created"
	EndpointDeleted       PubSubChannel = "endpoint:deleted"
	OrganizationModified  PubSubChannel = "organization:modified"
	SubscriptionModified  PubSubChannel = "subscription:modified"
	PricingModified       PubSubChannel = "pricing:modified"
	OrgPermissionModified PubSubChannel = "orgPermission:modified"
)

type RedisPrefix string

const (
	UsagePrefix        RedisPrefix = "usage"
	CountPrefix        RedisPrefix = "count"
	CostPrefix         RedisPrefix = "cost"
	TotalCostPrefix    RedisPrefix = "totalcost"
	TotalCountPrefix   RedisPrefix = "totalcount"
	PricingPrefix      RedisPrefix = "pricing"
	OrganizationPrefix RedisPrefix = "organization"
	EndpointPrefix     RedisPrefix = "endpoint"
	SubscriptionPrefix RedisPrefix = "subscription"
	PermissionPrefix   RedisPrefix = "permission"
)

var keyTypeTTLs = map[RedisPrefix]time.Duration{
	CountPrefix:        5 * time.Minute,
	CostPrefix:         5 * time.Minute,
	TotalCostPrefix:    60 * time.Minute,
	TotalCountPrefix:   60 * time.Minute,
	UsagePrefix:        24 * time.Hour,
	PricingPrefix:      24 * time.Hour,
	OrganizationPrefix: 24 * time.Hour,
	EndpointPrefix:     24 * time.Hour,
	SubscriptionPrefix: 24 * time.Hour,
}

func TTLFor(keyType RedisPrefix) time.Duration {
	if ttl, ok := keyTypeTTLs[keyType]; ok {
		return ttl
	}
	return time.Hour
}

func ParseUsageKey(key string) (orgID, subID, endpointID, timestamp int32, err error) {

	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("invalid key format")
	}
	org64, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid orgID")
	}
	sub64, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid subscriptionID")
	}
	endpoint64, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid endpointID")
	}
	timestamp64, err := strconv.ParseInt(parts[3], 10, 32)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid endpointID")
	}
	return int32(org64), int32(sub64), int32(endpoint64), int32(timestamp64), nil
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

func NextIntervalUnix(t time.Time, interval time.Duration) int64 {
	remainder := t.UnixNano() % interval.Nanoseconds()
	if remainder == 0 {
		return t.Unix()
	}
	return t.Add(time.Duration(interval.Nanoseconds() - remainder)).Unix()
}

type EndpointCreatedEvent struct {
	Code   string
	Path   string
	Method string
}

type EndpointDeletedEvent struct {
	Code string
}

type OrganizationModifiedEvent struct {
	ID   int32
	Name string
}

type OrgPermissionModifiedEvent struct {
	ID int32
}

type SubscriptionModifiedEvent struct {
	ID int32
}

type PricingModifiedEvent struct {
	ID   int32
	Type string
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

func ConvertProtoStruct(v interface{}) (*structpb.Struct, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var generic map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &generic); err != nil {
		return nil, err
	}

	return structpb.NewStruct(generic)
}
