package router

import (
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/handler"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/server"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OrgTypeHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/orgType")
	routerGrp.POST("/", h.CreateOrgTypeHandler)
	routerGrp.POST("/batch", h.CreateOrgTypeInBatchHandler)
	routerGrp.GET("/", h.ListOrgTypeHandler)
	routerGrp.DELETE("/:Id", h.DeleteOrgTypeHandler)
}

func SubTierHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/subTier")
	routerGrp.POST("/", h.CreateSubscriptionTierHandler)
	routerGrp.POST("/batch", h.CreateSubscriptionTierInBatchHandler)
	routerGrp.GET("/", h.ListSubscriptionTiersHandler)
	routerGrp.DELETE("/:Id", h.DeleteSubscriptionTierHandler)
}

func EndpointHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/endpoint")
	routerGrp.POST("/", h.RegisterEndpointHandler)
	routerGrp.POST("/batch", h.RegisterEndpointInBatchHandler)
	routerGrp.GET("/", h.ListEndpointsHandler)
	routerGrp.DELETE("/:Id", h.DeleteEndpointsByIdHandler)
}

func OrganizationHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/org")
	routerGrp.POST("/", h.CreateOrganizationandler)
	routerGrp.POST("/batch", h.CreateOrganizationInBatchandler)
	routerGrp.GET("/", h.ListOrganizationsHandler)
	routerGrp.DELETE("/:Id", h.DeleteOrganizationByIdHandler)
	routerGrp.GET("/:Id", h.GetOrganizationByIdHandler)
}

func TierPricingHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/tierPricing")
	routerGrp.POST("/batch", h.CreateTierPricingInBatchandler)
	routerGrp.DELETE("/tierId/:tier_id", h.DeleteTierPricingHandler)
	routerGrp.DELETE("/id/:id", h.DeleteTierPricingHandler)
	routerGrp.GET("/:tier_id", h.GetTierPricingByTierIdHandler)
}

func SubscriptionHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/subscription")
	routerGrp.POST("/", h.CreateSubscriptionHandler)
	routerGrp.POST("/batch", h.CreateSubscriptionInBatchandler)
	routerGrp.DELETE("/id/:id", h.DeleteSubscriptionHandler)
	routerGrp.DELETE("/orgId/:organization_id", h.DeleteSubscriptionHandler)
	routerGrp.GET("/id/:id", h.GetSubscriptionHandler)
	routerGrp.GET("/orgId/:organization_id", h.GetSubscriptionByrgIdHandler)
	routerGrp.GET("/", h.ListSubscriptionHandler)
}

func CustomPricingHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/customPricing")
	routerGrp.POST("/", h.CreateCustomPricingHandler)
	routerGrp.POST("/batch", h.CreateCustomPricingInBatchandler)
	routerGrp.DELETE("/subId/:subscription_id", h.DeleteCustomPricingHandler)
	routerGrp.DELETE("/id/:id", h.DeleteCustomPricingHandler)
	routerGrp.GET("/:subscription_id", h.GetCustomPricingHandler)
}

func ResourceTypeHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/resourceType")
	routerGrp.POST("/", h.CreateResurceTypeHandler)
	routerGrp.POST("/batch", h.CreateResurceTypeInBatchHandler)
	routerGrp.DELETE("/:id", h.DeleteResourceTypeHandler)
	routerGrp.GET("/", h.ListResourceTypeHandler)
}

func OrgPermissionHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/orgPermission")
	routerGrp.POST("/", h.CreateOrgPermissionHandler)
	routerGrp.POST("/batch", h.CreateOrgPermissionInBatchHandler)
	routerGrp.DELETE("/:organization_id", h.DeleteOrgPermissionHandler)
	routerGrp.GET("/:organization_id", h.GetOrgPermissionHandler)
}

func BillingHistoryHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/billingHistory")
	routerGrp.POST("/", h.CreateBillingHistoryHandler)
	routerGrp.POST("/batch", h.CreateBillingHistoryInBatchHandler)
	routerGrp.GET("/:id", h.GetBillingHistoryByIdHandler)
	routerGrp.GET("/orgId/:organization_id", h.GetBillingHistoryByOrgIdHandler)
	routerGrp.GET("/subId/:subscription_id", h.GetBillingHistoryBySubIdHandler)
}

func ApiUsageSummaryHandler(r *gin.RouterGroup, h *handler.AdminHandler) {
	routerGrp := r.Group("/apiUsageSummary")
	routerGrp.POST("/batch", h.CreateApiUsageInBatchHandler)
	routerGrp.GET("/orgId/:organization_id", h.GetApiUsageSummaryByOrgIdHandler)
	routerGrp.GET("/subId/:subscription_id", h.GetApiUsageSummaryBySubIdHandler)
	routerGrp.GET("/endpointId/:endpoint_id", h.GetApiUsageSummaryByEndpointIdHandler)
}

func RegisterHandlers(
	mainRouter *gin.RouterGroup,
	logger api.Logger,
	rw *server.ResponseWriter,
	db *sqlcgen.Queries,
	conn *pgxpool.Pool,
	validator *validator.Validate,
) {

	handler := handler.NewAdminHandler(logger, rw, db, conn, validator)

	mainRouter.GET("/", handler.RootHandler)

	OrgTypeHandler(mainRouter, handler)
	SubTierHandler(mainRouter, handler)
	EndpointHandler(mainRouter, handler)
	OrganizationHandler(mainRouter, handler)
	TierPricingHandler(mainRouter, handler)
	SubscriptionHandler(mainRouter, handler)
	CustomPricingHandler(mainRouter, handler)
	ResourceTypeHandler(mainRouter, handler)
	OrgPermissionHandler(mainRouter, handler)
	BillingHistoryHandler(mainRouter, handler)
	ApiUsageSummaryHandler(mainRouter, handler)
}
