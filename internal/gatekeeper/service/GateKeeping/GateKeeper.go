package gatekeeping

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
)

func (s *GateKeepingService) ValidateRequest(ctx context.Context, input *ValidateRequestInput) (*ValidationRequestOutput, error) {
	endpointCode, found := s.Match.Match(input.Method, input.Path)
	if !found {
		return nil, server.NewError(
			server.ErrorNotFound, "no matching endpoint", nil,
		)
	}

	org, err := caching.GetFromCache(ctx, s.Cache, "org:"+input.OrganizationName, func() (sqlcgen.GetOrganizationByNameRow, error) {
		return s.DB.GetOrganizationByName(ctx, input.OrganizationName)
	})
	if err != nil {
		return nil, server.NewError(
			server.ErrorNotFound, "organization not found", err,
		)
	}

	endpoint, err := caching.GetFromCache(ctx, s.Cache, "endpoint:"+endpointCode, func() (sqlcgen.GetEndpointByNameRow, error) {
		return s.DB.GetEndpointByName(ctx, endpointCode)
	})
	if err != nil {
		return nil, server.NewError(
			server.ErrorNotFound, "endpoint not found", err,
		)
	}

	cacheKey := "subscription:" + strconv.Itoa(int(org.ID)) + ":" + strconv.Itoa(int(endpoint.ID))
	sub, err := caching.GetFromCache(ctx, s.Cache, cacheKey, func() (sqlcgen.GetActiveSubscriptionRow, error) {
		return s.DB.GetActiveSubscription(ctx, sqlcgen.GetActiveSubscriptionParams{
			OrganizationID: org.ID,
			ApiEndpointID:  endpoint.ID,
		})
	})
	if err != nil || !sub.Active.Bool {
		return nil, server.NewError(
			server.ErrorUnauthorized, "no active subscription", nil,
		)
	}
	if sub.ExpiryTimestamp.Int32 > 0 && time.Now().Unix() > int64(sub.ExpiryTimestamp.Int32) {
		return nil, server.NewError(
			server.ErrorUnauthorized, "subscription expired", nil,
		)
	}

	var remaining *int32
	if sub.ApiLimit.Int32 > 0 {
		key := usageKey(org.ID, endpoint.ID)
		snapshot := s.Cache.GetAllLocalValues("usage")
		count := int32(snapshot[key])
		if count >= sub.ApiLimit.Int32 {
			return nil, server.NewError(
				server.ErrorUnauthorized, "quota exceeded", nil,
			)
		}
		left := sub.ApiLimit.Int32 - count
		remaining = &left
	}

	return &ValidationRequestOutput{
		Organization: org,
		Endpoint:     endpoint,
		Subscription: sub,
		Remaining:    remaining,
	}, nil
}

func (s *GateKeepingService) RecordUsage(ctx context.Context, input *RecordUsageInput) (float64, error) {

	endpointCode, found := s.Match.Match(input.Method, input.Path)
	if !found {
		return 0, server.NewError(
			server.ErrorInternal, "no matching endpoint", nil,
		)
	}

	org, err := s.DB.GetOrganizationByName(ctx, input.OrganizationName)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal, "organization not found", err,
		)
	}

	endpoint, err := s.DB.GetEndpointByName(ctx, endpointCode)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal, "endpoint not found", err,
		)
	}

	sub, err := s.DB.GetActiveSubscription(ctx, sqlcgen.GetActiveSubscriptionParams{
		OrganizationID: org.ID,
		ApiEndpointID:  endpoint.ID,
	})
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal, "subscription not found", err,
		)
	}

	pricing, err := s.DB.GetPricing(ctx, sqlcgen.GetPricingParams{
		SubscriptionID: sub.ID,
		ApiEndpointID:  endpoint.ID,
	})
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal, "pricing error", err,
		)
	}

	s.Cache.IncrementLocalValue("usage", usageKey(org.ID, endpoint.ID), 1)
	return pricing, nil
}

func usageKey(orgID, endpointID int32) string {
	return fmt.Sprintf("%d:%d", orgID, endpointID)
}
