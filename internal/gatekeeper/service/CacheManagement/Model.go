package cachemanagement

import (
	"fmt"
	"strconv"
	"strings"
)

// Utility to parse keys like: usage:{orgID}:{endpointID}
func ParseUsageKey(key string) (orgID, endpointID int32, err error) {
	parts := strings.Split(key, ":")
	if len(parts) != 3 {
		return 0, 0, fmt.Errorf("invalid key format")
	}
	org64, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid orgID")
	}
	endpoint64, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid endpointID")
	}
	return int32(org64), int32(endpoint64), nil
}
