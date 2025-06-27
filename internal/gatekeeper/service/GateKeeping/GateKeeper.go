package gatekeeping

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
)

func (s *GateKeepingService) ValidateRequest(ctx context.Context, orgName, endpointName string) error {

	orgAny, err := s.Cache.Get(ctx, "org:"+orgName, func() (interface{}, error) {
		return s.DB.GetOrganizationByName(ctx, orgName)
	})
	if err != nil {
		return errors.New("organization not found")
	}

	endpointAny, err := s.Cache.Get(ctx, "endpoint:"+endpointName, func() (interface{}, error) {
		return s.DB.GetEndpointByName(ctx, endpointName)
	})
	if err != nil {
		return errors.New("endpoint not found")
	}

	org := orgAny.(Organization)
	endpoint := endpointAny.(ApiEndpoint)
	sub, err := s.DB.GetActiveSubscription(ctx, sqlcgen.GetActiveSubscriptionParams{
		OrganizationID: org.ID,
		ApiEndpointID:  endpoint.ID,
	})
	if err != nil || !sub.Active.Bool {
		return errors.New("no active subscription")
	}
	if sub.ExpiryTimestamp.Int32 > 0 && time.Now().Unix() > int64(sub.ExpiryTimestamp.Int32) {
		return errors.New("subscription expired")
	}

	// Usage check from cache
	key := usageCacheKey(org.ID, endpoint.ID)
	val, _ := s.Cache.Get(ctx, key, func() (interface{}, error) {
		return 0, nil
	})
	if count, ok := val.(int); ok && int32(count) >= sub.ApiLimit.Int32 && sub.ApiLimit.Int32 > 0 {
		return errors.New("quota exceeded")
	}

	return nil
}

func (s *GateKeepingService) RecordUsage(ctx context.Context, orgName, endpointName string) (float64, error) {
	org, err := s.DB.GetOrganizationByName(ctx, orgName)
	if err != nil {
		return 0, errors.New("organization not found")
	}

	endpoint, err := s.DB.GetEndpointByName(ctx, endpointName)
	if err != nil {
		return 0, errors.New("endpoint not found")
	}

	sub, err := s.DB.GetActiveSubscription(ctx, sqlcgen.GetActiveSubscriptionParams{
		OrganizationID: org.ID,
		ApiEndpointID:  endpoint.ID,
	})
	if err != nil {
		return 0, errors.New("subscription not found")
	}

	pricing, err := s.DB.GetPricing(ctx, sqlcgen.GetPricingParams{
		SubscriptionID: sub.ID,
		ApiEndpointID:  endpoint.ID,
	})
	if err != nil {
		return 0, errors.New("pricing error")
	}

	key := usageCacheKey(org.ID, endpoint.ID)

	// Try to get current usage count from cache
	val, _ := s.Cache.Get(ctx, key, func() (interface{}, error) {
		return 0, nil // Default count if not in cache
	})

	if count, ok := val.(int); ok {
		_ = s.Cache.Set(ctx, key, count+1)
	} else {
		_ = s.Cache.Set(ctx, key, 1)
	}

	// Async DB usage record
	go s.DB.IncrementUsage(context.Background(), sqlcgen.IncrementUsageParams{
		OrganizationID: org.ID,
		SubscriptionID: sub.ID,
		ApiEndpointID:  endpoint.ID,
	})

	return pricing, nil
}

func usageCacheKey(orgID, endpointID int32) string {
	return fmt.Sprintf("usage:%d:%d", orgID, endpointID)
}
