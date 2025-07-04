package gatekeeping

import (
	"context"
	"strconv"
	"time"

	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/counter"
	"github.com/bignyap/go-utilities/server"
)

func (s *GateKeepingService) ValidateRequest(ctx context.Context, input *ValidateRequestInput) (*ValidationRequestOutput, error) {

	orgSubDetails, err := s.GetOrgSubDetailsFromCache(ctx, input.Method, input.Path, input.OrganizationName)
	if err != nil {
		return nil, server.NewError(
			server.ErrorUnauthorized, "failed to validate request", err,
		)
	}

	return &orgSubDetails.ValidationRequestOutput, nil
}

func (s *GateKeepingService) RecordUsage(ctx context.Context, input *RecordUsageInput) (float64, error) {

	orgSubDetails, err := s.GetOrgSubDetailsFromCache(ctx, input.Method, input.Path, input.OrganizationName)
	if err != nil {
		return 0.0, server.NewError(
			server.ErrorUnauthorized, "failed to validate request", err,
		)
	}

	pricingcacheKey := common.RedisKeyFormatter(
		string(common.PricingKeyPrefix), orgSubDetails.Organization.Name, orgSubDetails.EndpointCode,
	)
	pricing, err := caching.GetFromCache(ctx, s.Cache, pricingcacheKey, func() (float64, error) {
		return s.DB.GetPricing(ctx, sqlcgen.GetPricingParams{
			SubscriptionID: orgSubDetails.Subscription.ID,
			ApiEndpointID:  orgSubDetails.Endpoint.ID,
		})
	})
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal, "pricing error", err,
		)
	}

	updateUsageCounters(s.CounterWorker, orgSubDetails, pricing)

	return pricing, nil
}

func updateUsageCounters(countWorker *counter.CounterWorker, orgSubDetails *GetOrgSubDetailsOutput, pricing float64) {

	usageKey := common.UsageKey(
		orgSubDetails.Organization.ID,
		orgSubDetails.Subscription.ID,
		orgSubDetails.Endpoint.ID,
	)
	countWorker.Increment(
		string(common.Usageprefix),
		common.RedisKeyFormatter(usageKey, string(common.Costprefix)),
		pricing,
	)
	countWorker.Increment(
		string(common.Usageprefix),
		common.RedisKeyFormatter(usageKey, string(common.Countprefix)),
		1,
	)

	totalUsageKey := common.RedisKeyFormatter(
		orgSubDetails.Organization.ID,
		orgSubDetails.Subscription.ID,
	)
	countWorker.Increment(
		string(common.Usageprefix),
		common.RedisKeyFormatter(totalUsageKey, string(common.TotalCostKeyPrefix)),
		pricing,
	)
	countWorker.Increment(
		string(common.Usageprefix),
		common.RedisKeyFormatter(totalUsageKey, string(common.TotalCountKeyPrefix)),
		1,
	)
}

func (s *GateKeepingService) GetOrgSubDetailsFromCache(ctx context.Context, method, path, orgName string) (*GetOrgSubDetailsOutput, error) {

	endpointCode, found := s.Match.Match(method, path)
	if !found {
		return nil, server.NewError(
			server.ErrorNotFound, "no matching endpoint", nil,
		)
	}

	orgKey := common.RedisKeyFormatter(string(common.OrganizationKeyPrefix), orgName)
	org, err := caching.GetFromCache(ctx, s.Cache, orgKey, func() (sqlcgen.GetOrganizationByNameRow, error) {
		return s.DB.GetOrganizationByName(ctx, orgName)
	})
	if err != nil {
		return nil, server.NewError(
			server.ErrorNotFound, "organization not found", err,
		)
	}

	epKey := common.RedisKeyFormatter(string(common.EndpointKeyPrefix), endpointCode)
	endpoint, err := caching.GetFromCache(ctx, s.Cache, epKey, func() (sqlcgen.GetEndpointByNameRow, error) {
		return s.DB.GetEndpointByName(ctx, endpointCode)
	})
	if err != nil {
		return nil, server.NewError(
			server.ErrorNotFound, "endpoint not found", err,
		)
	}

	subKey := common.RedisKeyFormatter(
		string(common.SubscriptionKeyPrefix),
		strconv.Itoa(int(org.ID)),
		strconv.Itoa(int(endpoint.ID)),
	)
	sub, err := caching.GetFromCache(ctx, s.Cache, subKey, func() (sqlcgen.GetActiveSubscriptionRow, error) {
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

		usage, err := s.GetUsageDetailFromCache(ctx, org.ID, sub.ID, endpoint.ID)
		if err != nil {
			return nil, server.NewError(
				server.ErrorInternal, "error fetching the total usage", err,
			)
		}
		if usage >= sub.ApiLimit.Int32 {
			return nil, server.NewError(
				server.ErrorUnauthorized, "quota exceeded", nil,
			)
		}
		left := sub.ApiLimit.Int32 - usage
		remaining = &left
	}

	return &GetOrgSubDetailsOutput{
		ValidationRequestOutput: ValidationRequestOutput{
			Organization: org,
			Endpoint:     endpoint,
			Subscription: sub,
			Remaining:    remaining,
		},
		EndpointCode: endpointCode,
	}, nil
}

func (s *GateKeepingService) GetUsageDetailFromCache(ctx context.Context, orgId, subId, endpointId int32) (int32, error) {

	usageRedisKey := common.RedisKeyFormatter(
		string(common.Usageprefix), string(common.TotalCostKeyPrefix),
		common.UsageKey(orgId, subId, endpointId),
	)
	value, err := s.Cache.Get(ctx, usageRedisKey, func() (interface{}, error) {
		// Fetch from DB if not found in cache
		orgUsageDetails, err := s.DB.GetQuotaUsageBySubscriptionID(ctx, subId)
		if err != nil {
			return 0, server.NewError(
				server.ErrorInternal, "error fetching usage details", err,
			)
		}

		if orgUsageDetails.CallsUsed == 0 {
			return 0, nil // No usage found
		}
		return int32(orgUsageDetails.CallsUsed), nil
	})
	if err != nil {
		return 0, err
	}

	// Ensure the value is cast to int32
	if usage, ok := value.(int32); ok {
		return usage, nil
	}

	return 0, server.NewError(
		server.ErrorInternal, "unexpected value type in cache", nil,
	)
}
