package gatekeeping

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/counter"
	"github.com/bignyap/go-utilities/server"
)

func (s *GateKeepingService) FlushAllCache(ctx context.Context) {
	s.CounterWorker.FlushNow("*", ctx)
	s.Cache.ResetRedisValues(ctx, "*")
}

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
		string(common.PricingPrefix), orgSubDetails.Organization.Name, orgSubDetails.EndpointCode,
	)

	pricing, err := caching.GetFromCache(ctx, s.Cache, pricingcacheKey, func() (sqlcgen.GetPricingRow, error) {
		val, err := s.DB.GetPricing(ctx, sqlcgen.GetPricingParams{
			SubscriptionID: orgSubDetails.Subscription.ID,
			ApiEndpointID:  orgSubDetails.Endpoint.ApiEndpointID,
		})
		return val, err
	})
	if err != nil {
		return 0, server.NewError(server.ErrorInternal, "pricing error", err)
	}

	effectiveCalls := 1
	if val := ctx.Value("x-effective-calls"); val != nil {
		if intVal, ok := val.(int); ok && intVal > 0 {
			effectiveCalls = intVal
		}
	}
	effectivePricing := pricing.CostPerCall
	if strings.EqualFold(pricing.CostMode, "dynamic") {
		effectivePricing *= float64(effectiveCalls)
	}

	updateUsageCounters(s.CounterWorker, orgSubDetails, s.FlushInterval, effectivePricing)

	return effectivePricing, nil
}

func updateUsageCounters(
	countWorker *counter.CounterWorker,
	orgSubDetails *GetOrgSubDetailsOutput,
	interval int64,
	pricing float64,
) {

	timestamp := common.NextIntervalUnix(time.Now(), time.Duration(interval)*time.Second)
	timestampStr := strconv.FormatInt(timestamp, 10)

	usageKey := common.UsageKey(
		orgSubDetails.Organization.ID,
		orgSubDetails.Subscription.ID,
		orgSubDetails.Endpoint.ApiEndpointID,
	)

	countWorker.Increment(
		string(common.UsagePrefix),
		common.RedisKeyFormatter(usageKey, timestampStr, string(common.CostPrefix)),
		pricing,
	)

	countWorker.Increment(
		string(common.UsagePrefix),
		common.RedisKeyFormatter(usageKey, timestampStr, string(common.CountPrefix)),
		1,
	)

	totalUsageKey := common.RedisKeyFormatter(
		orgSubDetails.Organization.ID,
		orgSubDetails.Subscription.ID,
	)

	countWorker.Increment(
		string(common.UsagePrefix),
		common.RedisKeyFormatter(totalUsageKey, timestampStr, string(common.TotalCostPrefix)),
		pricing,
	)

	countWorker.Increment(
		string(common.UsagePrefix),
		common.RedisKeyFormatter(totalUsageKey, timestampStr, string(common.TotalCountPrefix)),
		1,
	)
}

func (s *GateKeepingService) GetOrgSubDetailsFromCache(ctx context.Context, method, path, orgName string) (*GetOrgSubDetailsOutput, error) {
	endpointCode, found := s.Match.Match(method, path)
	if !found {
		return nil, server.NewError(server.ErrorNotFound, "no matching endpoint", nil)
	}

	orgKey := common.RedisKeyFormatter(string(common.OrganizationPrefix), orgName)
	org, err := caching.GetFromCache(ctx, s.Cache, orgKey, func() (sqlcgen.GetOrganizationByNameRow, error) {
		return s.DB.GetOrganizationByName(ctx, orgName)
	})
	if err != nil {
		return nil, server.NewError(server.ErrorNotFound, "organization not found", err)
	}

	epKey := common.RedisKeyFormatter(string(common.EndpointPrefix), endpointCode)
	endpoint, err := caching.GetFromCache(ctx, s.Cache, epKey, func() (sqlcgen.GetApiEndpointByNameRow, error) {
		return s.DB.GetApiEndpointByName(ctx, endpointCode)
	})
	if err != nil {
		return nil, server.NewError(server.ErrorNotFound, "endpoint not found", err)
	}

	epPerKey := common.RedisKeyFormatter(
		string(common.EndpointPrefix), string(endpoint.ResourceTypeID),
		string(common.PermissionPrefix), endpoint.PermissionCode,
	)
	orgPerExists, err := caching.GetFromCache(ctx, s.Cache, epPerKey, func() (bool, error) {
		return s.DB.CheckOrgPermission(ctx, sqlcgen.CheckOrgPermissionParams{
			ResourceTypeID: endpoint.ResourceTypeID,
			PermissionCode: endpoint.PermissionCode,
		})
	})
	if err != nil || !orgPerExists {
		return nil, server.NewError(server.ErrorUnauthorized, "insufficient permission", err)
	}

	subKey := common.RedisKeyFormatter(
		string(common.SubscriptionPrefix),
		strconv.Itoa(int(org.ID)),
		strconv.Itoa(int(endpoint.ApiEndpointID)),
	)
	sub, err := caching.GetFromCache(ctx, s.Cache, subKey, func() (sqlcgen.GetActiveSubscriptionRow, error) {
		return s.DB.GetActiveSubscription(ctx, org.ID)
	})
	if err != nil || !sub.Active.Bool {
		return nil, server.NewError(server.ErrorUnauthorized, "no active subscription", nil)
	}
	if sub.ExpiryTimestamp.Int32 > 0 && time.Now().Unix() > int64(sub.ExpiryTimestamp.Int32) {
		return nil, server.NewError(server.ErrorUnauthorized, "subscription expired", nil)
	}

	// Check if the ratelimit had reached for the endpoint + subscriptionid
	// We might need to implement a ratelimiter for this. Let's think through this
	// I definitely do not want to have a full blown Ratelimiter implementation here

	var remaining int32 = -1
	if sub.ApiLimit.Int32 > 0 {
		usage, err := s.GetUsageDetailFromCache(ctx, org.ID, sub.ID, endpoint.ApiEndpointID)
		if err != nil {
			return nil, server.NewError(server.ErrorInternal, "error fetching the total usage", err)
		}
		if usage >= sub.ApiLimit.Int32 {
			return nil, server.NewError(server.ErrorUnauthorized, "quota exceeded", nil)
		}
		remaining = max(sub.ApiLimit.Int32-usage, 0)
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

	timestamp := common.NextIntervalUnix(time.Now(), time.Duration(s.FlushInterval)*time.Second)
	timestampStr := strconv.FormatInt(timestamp, 10)

	usageRedisKey := common.RedisKeyFormatter(
		string(common.UsagePrefix), string(orgId), string(subId),
		timestampStr, string(common.TotalCostPrefix),
	)

	usage, err := caching.GetFromCache(ctx, s.Cache, usageRedisKey, func() (int32, error) {
		orgUsageDetails, err := s.DB.GetQuotaUsageBySubscriptionID(ctx, subId)
		if err != nil {
			return 0, server.NewError(server.ErrorInternal, "error fetching usage details", err)
		}
		return orgUsageDetails.CostsUsed, nil
	})
	if err != nil {
		return 0, err
	}
	return usage, nil
}
