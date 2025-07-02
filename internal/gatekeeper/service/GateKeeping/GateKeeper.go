package gatekeeping

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
)

func (s *GateKeepingService) ValidateRequest(ctx context.Context, input *ValidateRequestInput) error {

	endpointCode, found := s.Match.Match(input.Method, input.Path)
	if !found {
		return errors.New("no matching endpoint")
	}

	orgAny, err := s.Cache.Get(ctx, "org:"+input.OrganizationName, func() (interface{}, error) {
		return s.DB.GetOrganizationByName(ctx, input.OrganizationName)
	})
	if err != nil {
		return errors.New("organization not found")
	}

	endpointAny, err := s.Cache.Get(ctx, "endpoint:"+endpointCode, func() (interface{}, error) {
		return s.DB.GetEndpointByName(ctx, endpointCode)
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

	if sub.ApiLimit.Int32 > 0 {
		key := usageKey(org.ID, endpoint.ID)
		snapshot := s.Cache.GetAllLocalValues("usage")
		if count := snapshot[key]; int32(count) >= sub.ApiLimit.Int32 {
			return errors.New("quota exceeded")
		}
	}

	return nil
}

func (s *GateKeepingService) RecordUsage(ctx context.Context, input *RecordUsageInput) (float64, error) {
	endpointCode, found := s.Match.Match(input.Method, input.Path)
	if !found {
		return 0, errors.New("no matching endpoint")
	}

	org, err := s.DB.GetOrganizationByName(ctx, input.OrganizationName)
	if err != nil {
		return 0, errors.New("organization not found")
	}

	endpoint, err := s.DB.GetEndpointByName(ctx, endpointCode)
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

	s.Cache.IncrementLocalValue("usage", usageKey(org.ID, endpoint.ID), 1)
	return pricing, nil
}

func usageKey(orgID, endpointID int32) string {
	return fmt.Sprintf("%d:%d", orgID, endpointID)
}
