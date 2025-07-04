// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: copyfrom.go

package sqlcgen

import (
	"context"
)

// iteratorForCreateApiUsageSummaries implements pgx.CopyFromSource.
type iteratorForCreateApiUsageSummaries struct {
	rows                 []CreateApiUsageSummariesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateApiUsageSummaries) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateApiUsageSummaries) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].UsageStartDate,
		r.rows[0].UsageEndDate,
		r.rows[0].TotalCalls,
		r.rows[0].TotalCost,
		r.rows[0].SubscriptionID,
		r.rows[0].ApiEndpointID,
		r.rows[0].OrganizationID,
	}, nil
}

func (r iteratorForCreateApiUsageSummaries) Err() error {
	return nil
}

func (q *Queries) CreateApiUsageSummaries(ctx context.Context, arg []CreateApiUsageSummariesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"api_usage_summary"}, []string{"usage_start_date", "usage_end_date", "total_calls", "total_cost", "subscription_id", "api_endpoint_id", "organization_id"}, &iteratorForCreateApiUsageSummaries{rows: arg})
}

// iteratorForCreateBillingHistories implements pgx.CopyFromSource.
type iteratorForCreateBillingHistories struct {
	rows                 []CreateBillingHistoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateBillingHistories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateBillingHistories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].BillingStartDate,
		r.rows[0].BillingEndDate,
		r.rows[0].TotalAmountDue,
		r.rows[0].TotalCalls,
		r.rows[0].PaymentStatus,
		r.rows[0].PaymentDate,
		r.rows[0].BillingCreatedAt,
		r.rows[0].SubscriptionID,
	}, nil
}

func (r iteratorForCreateBillingHistories) Err() error {
	return nil
}

func (q *Queries) CreateBillingHistories(ctx context.Context, arg []CreateBillingHistoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"billing_history"}, []string{"billing_start_date", "billing_end_date", "total_amount_due", "total_calls", "payment_status", "payment_date", "billing_created_at", "subscription_id"}, &iteratorForCreateBillingHistories{rows: arg})
}

// iteratorForCreateCustomPricings implements pgx.CopyFromSource.
type iteratorForCreateCustomPricings struct {
	rows                 []CreateCustomPricingsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateCustomPricings) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateCustomPricings) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].CustomCostPerCall,
		r.rows[0].CustomRateLimit,
		r.rows[0].SubscriptionID,
		r.rows[0].TierBasePricingID,
	}, nil
}

func (r iteratorForCreateCustomPricings) Err() error {
	return nil
}

func (q *Queries) CreateCustomPricings(ctx context.Context, arg []CreateCustomPricingsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"custom_endpoint_pricing"}, []string{"custom_cost_per_call", "custom_rate_limit", "subscription_id", "tier_base_pricing_id"}, &iteratorForCreateCustomPricings{rows: arg})
}

// iteratorForCreateOrgPermissions implements pgx.CopyFromSource.
type iteratorForCreateOrgPermissions struct {
	rows                 []CreateOrgPermissionsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateOrgPermissions) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateOrgPermissions) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ResourceTypeID,
		r.rows[0].PermissionCode,
		r.rows[0].OrganizationID,
	}, nil
}

func (r iteratorForCreateOrgPermissions) Err() error {
	return nil
}

func (q *Queries) CreateOrgPermissions(ctx context.Context, arg []CreateOrgPermissionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"organization_permission"}, []string{"resource_type_id", "permission_code", "organization_id"}, &iteratorForCreateOrgPermissions{rows: arg})
}

// iteratorForCreateOrgTypes implements pgx.CopyFromSource.
type iteratorForCreateOrgTypes struct {
	rows                 []string
	skippedFirstNextCall bool
}

func (r *iteratorForCreateOrgTypes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateOrgTypes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0],
	}, nil
}

func (r iteratorForCreateOrgTypes) Err() error {
	return nil
}

func (q *Queries) CreateOrgTypes(ctx context.Context, organizationTypeName []string) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"organization_type"}, []string{"organization_type_name"}, &iteratorForCreateOrgTypes{rows: organizationTypeName})
}

// iteratorForCreateOrganizations implements pgx.CopyFromSource.
type iteratorForCreateOrganizations struct {
	rows                 []CreateOrganizationsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateOrganizations) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateOrganizations) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].OrganizationName,
		r.rows[0].OrganizationCreatedAt,
		r.rows[0].OrganizationUpdatedAt,
		r.rows[0].OrganizationRealm,
		r.rows[0].OrganizationCountry,
		r.rows[0].OrganizationSupportEmail,
		r.rows[0].OrganizationActive,
		r.rows[0].OrganizationReportQ,
		r.rows[0].OrganizationConfig,
		r.rows[0].OrganizationTypeID,
	}, nil
}

func (r iteratorForCreateOrganizations) Err() error {
	return nil
}

func (q *Queries) CreateOrganizations(ctx context.Context, arg []CreateOrganizationsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"organization"}, []string{"organization_name", "organization_created_at", "organization_updated_at", "organization_realm", "organization_country", "organization_support_email", "organization_active", "organization_report_q", "organization_config", "organization_type_id"}, &iteratorForCreateOrganizations{rows: arg})
}

// iteratorForCreateResourceTypes implements pgx.CopyFromSource.
type iteratorForCreateResourceTypes struct {
	rows                 []CreateResourceTypesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateResourceTypes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateResourceTypes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ResourceTypeName,
		r.rows[0].ResourceTypeCode,
		r.rows[0].ResourceTypeDescription,
	}, nil
}

func (r iteratorForCreateResourceTypes) Err() error {
	return nil
}

func (q *Queries) CreateResourceTypes(ctx context.Context, arg []CreateResourceTypesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"resource_type"}, []string{"resource_type_name", "resource_type_code", "resource_type_description"}, &iteratorForCreateResourceTypes{rows: arg})
}

// iteratorForCreateSubscriptionTiers implements pgx.CopyFromSource.
type iteratorForCreateSubscriptionTiers struct {
	rows                 []CreateSubscriptionTiersParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateSubscriptionTiers) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateSubscriptionTiers) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].TierName,
		r.rows[0].TierDescription,
		r.rows[0].TierCreatedAt,
		r.rows[0].TierUpdatedAt,
	}, nil
}

func (r iteratorForCreateSubscriptionTiers) Err() error {
	return nil
}

func (q *Queries) CreateSubscriptionTiers(ctx context.Context, arg []CreateSubscriptionTiersParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"subscription_tier"}, []string{"tier_name", "tier_description", "tier_created_at", "tier_updated_at"}, &iteratorForCreateSubscriptionTiers{rows: arg})
}

// iteratorForCreateSubscriptions implements pgx.CopyFromSource.
type iteratorForCreateSubscriptions struct {
	rows                 []CreateSubscriptionsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateSubscriptions) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateSubscriptions) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].SubscriptionName,
		r.rows[0].SubscriptionType,
		r.rows[0].SubscriptionCreatedDate,
		r.rows[0].SubscriptionUpdatedDate,
		r.rows[0].SubscriptionStartDate,
		r.rows[0].SubscriptionApiLimit,
		r.rows[0].SubscriptionExpiryDate,
		r.rows[0].SubscriptionDescription,
		r.rows[0].SubscriptionStatus,
		r.rows[0].OrganizationID,
		r.rows[0].SubscriptionTierID,
		r.rows[0].SubscriptionBillingInterval,
		r.rows[0].SubscriptionBillingModel,
		r.rows[0].SubscriptionQuotaResetInterval,
	}, nil
}

func (r iteratorForCreateSubscriptions) Err() error {
	return nil
}

func (q *Queries) CreateSubscriptions(ctx context.Context, arg []CreateSubscriptionsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"subscription"}, []string{"subscription_name", "subscription_type", "subscription_created_date", "subscription_updated_date", "subscription_start_date", "subscription_api_limit", "subscription_expiry_date", "subscription_description", "subscription_status", "organization_id", "subscription_tier_id", "subscription_billing_interval", "subscription_billing_model", "subscription_quota_reset_interval"}, &iteratorForCreateSubscriptions{rows: arg})
}

// iteratorForCreateTierPricings implements pgx.CopyFromSource.
type iteratorForCreateTierPricings struct {
	rows                 []CreateTierPricingsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateTierPricings) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateTierPricings) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].SubscriptionTierID,
		r.rows[0].ApiEndpointID,
		r.rows[0].BaseCostPerCall,
		r.rows[0].BaseRateLimit,
	}, nil
}

func (r iteratorForCreateTierPricings) Err() error {
	return nil
}

func (q *Queries) CreateTierPricings(ctx context.Context, arg []CreateTierPricingsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"tier_base_pricing"}, []string{"subscription_tier_id", "api_endpoint_id", "base_cost_per_call", "base_rate_limit"}, &iteratorForCreateTierPricings{rows: arg})
}

// iteratorForRegisterApiEndpoints implements pgx.CopyFromSource.
type iteratorForRegisterApiEndpoints struct {
	rows                 []RegisterApiEndpointsParams
	skippedFirstNextCall bool
}

func (r *iteratorForRegisterApiEndpoints) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForRegisterApiEndpoints) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].EndpointName,
		r.rows[0].EndpointDescription,
		r.rows[0].HttpMethod,
		r.rows[0].PathTemplate,
		r.rows[0].ResourceTypeID,
	}, nil
}

func (r iteratorForRegisterApiEndpoints) Err() error {
	return nil
}

func (q *Queries) RegisterApiEndpoints(ctx context.Context, arg []RegisterApiEndpointsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"api_endpoint"}, []string{"endpoint_name", "endpoint_description", "http_method", "path_template", "resource_type_id"}, &iteratorForRegisterApiEndpoints{rows: arg})
}
