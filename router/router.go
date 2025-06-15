package router

import (
	"net/http"

	"github.com/bignyap/go-admin/handler"
)

func OrgTypeHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /orgType", apiConfig.CreateOrgTypeHandler)
	mux.HandleFunc("POST /orgType/batch", apiConfig.CreateOrgTypeInBatchHandler)
	mux.HandleFunc("GET /orgType", apiConfig.ListOrgTypeHandler)
	mux.HandleFunc("DELETE /orgType/{Id}", apiConfig.DeleteOrgTypeHandler)

}

func SubTierHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /subTier", apiConfig.CreateSubcriptionTierHandler)
	mux.HandleFunc("POST /subTier/batch", apiConfig.CreateSubscriptionTierInBatchHandler)
	mux.HandleFunc("GET /subTier", apiConfig.ListSubscriptionTiersHandler)
	mux.HandleFunc("DELETE /subTier/{Id}", apiConfig.DeleteSubscriptionTierHandler)

}

func EndpointHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /endpoint", apiConfig.RegisterEndpointHandler)
	mux.HandleFunc("POST /endpoint/batch", apiConfig.RegisterEndpointInBatchHandler)
	mux.HandleFunc("GET /endpoint", apiConfig.ListEndpointsHandler)
	mux.HandleFunc("DELETE /endpoint/{Id}", apiConfig.DeleteEndpointsByIdHandler)

}

func OrganizationHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /org", apiConfig.CreateOrganizationandler)
	mux.HandleFunc("POST /org/batch", apiConfig.CreateOrganizationInBatchandler)
	mux.HandleFunc("GET /org", apiConfig.ListOrganizationsHandler)
	mux.HandleFunc("DELETE /org/{Id}", apiConfig.DeleteOrganizationByIdHandler)
	mux.HandleFunc("GET /org/{Id}", apiConfig.GetOrganizationByIdHandler)

}

func TierPricingHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /tierPricing", apiConfig.CreateTierPricingHandler)
	mux.HandleFunc("POST /tierPricing/batch", apiConfig.CreateTierPricingInBatchandler)
	mux.HandleFunc("DELETE /tierPricing/tierId/{tier_id}", apiConfig.DeleteTierPricingHandler)
	mux.HandleFunc("DELETE /tierPricing/id/{id}", apiConfig.DeleteTierPricingHandler)
	mux.HandleFunc("GET /tierPricing/{tier_id}", apiConfig.GetTierPricingByTierIdHandler)

}

func SubscriptionHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /subscription", apiConfig.CreateSubscriptionHandler)
	mux.HandleFunc("POST /subscription/batch", apiConfig.CreateSubscriptionInBatchandler)
	mux.HandleFunc("DELETE /subscription/id/{id}", apiConfig.DeleteSubscriptionHandler)
	mux.HandleFunc("DELETE /subscription/orgId/{organization_id}", apiConfig.DeleteSubscriptionHandler)
	mux.HandleFunc("GET /subscription/id/{id}", apiConfig.GetSubscriptionHandler)
	mux.HandleFunc("GET /subscription/orgId/{organization_id}", apiConfig.GetSubscriptionByrgIdHandler)
	mux.HandleFunc("GET /subscription", apiConfig.ListSubscriptionHandler)

}

func CustomPricingHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /customPricing", apiConfig.CreateCustomPricingHandler)
	mux.HandleFunc("POST /customPricing/batch", apiConfig.CreateCustomPricingInBatchandler)
	mux.HandleFunc("DELETE /customPricing/subId/{subscription_id}", apiConfig.DeleteCustomPricingHandler)
	mux.HandleFunc("DELETE /customPricing/id/{id}", apiConfig.DeleteCustomPricingHandler)
	mux.HandleFunc("GET /customPricing/{subscription_id}", apiConfig.GetCustomPricingHandler)

}

func ResourceTypeHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /resourceType", apiConfig.CreateResurceTypeHandler)
	mux.HandleFunc("POST /resourceType/batch", apiConfig.CreateResurceTypeInBatchHandler)
	mux.HandleFunc("DELETE /resourceType/{id}", apiConfig.DeleteResourceTypeHandler)
	mux.HandleFunc("GET /resourceType", apiConfig.ListResourceTypeHandler)

}

func OrgPermissionHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /orgPermission", apiConfig.CreateOrgPermissionHandler)
	mux.HandleFunc("POST /orgPermission/batch", apiConfig.CreateOrgPermissionInBatchHandler)
	mux.HandleFunc("DELETE /orgPermission/{organization_id}", apiConfig.DeleteOrgPermissionHandler)
	mux.HandleFunc("GET /orgPermission/{organization_id}", apiConfig.GetOrgPermissionHandler)

}

func BillingHistoryHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /billingHistory", apiConfig.CreateBillingHistoryHandler)
	mux.HandleFunc("POST /billingHistory/batch", apiConfig.CreateBillingHistoryInBatchHandler)
	mux.HandleFunc("GET /billingHistory/{id}", apiConfig.GetBillingHistoryByIdHandler)
	mux.HandleFunc("GET /billingHistory/orgId/{organization_id}", apiConfig.GetBillingHistoryByOrgIdHandler)
	mux.HandleFunc("GET /billingHistory/subId/{subscription_id}", apiConfig.GetBillingHistoryBySubIdHandler)

}

func ApiUsageSummaryHandler(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("POST /apiUsageSummary", apiConfig.CreateApiUsageHander)
	mux.HandleFunc("POST /apiUsageSummary/batch", apiConfig.CreateApiUsageInBatchHander)
	mux.HandleFunc("GET /apiUsageSummary/orgId/{organization_id}", apiConfig.GetApiUsageSummaryByOrgIdHandler)
	mux.HandleFunc("GET /apiUsageSummary/subId/{subscription_id}", apiConfig.GetApiUsageSummaryBySubIdHandler)
	mux.HandleFunc("GET /apiUsageSummary/endpointId/{endpoint_id}", apiConfig.GetApiUsageSummaryByEndpointIdHandler)
}

func RegisterHandlers(mux *http.ServeMux, apiConfig *handler.ApiConfig) {

	mux.HandleFunc("/", handler.RootHandler)
	OrgTypeHandler(mux, apiConfig)
	SubTierHandler(mux, apiConfig)
	EndpointHandler(mux, apiConfig)
	OrganizationHandler(mux, apiConfig)
	TierPricingHandler(mux, apiConfig)
	SubscriptionHandler(mux, apiConfig)
	CustomPricingHandler(mux, apiConfig)
	ResourceTypeHandler(mux, apiConfig)
	OrgPermissionHandler(mux, apiConfig)
	BillingHistoryHandler(mux, apiConfig)
	ApiUsageSummaryHandler(mux, apiConfig)
}
