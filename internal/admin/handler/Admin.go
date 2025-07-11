package adminHandler

import (
	service "github.com/bignyap/go-admin/internal/admin/service/Billing"
	dashboard "github.com/bignyap/go-admin/internal/admin/service/Dashboard"
	organization "github.com/bignyap/go-admin/internal/admin/service/Organization"
	pricing "github.com/bignyap/go-admin/internal/admin/service/Pricing"
	resource "github.com/bignyap/go-admin/internal/admin/service/Resource"
	subscription "github.com/bignyap/go-admin/internal/admin/service/Subscription"
	usage "github.com/bignyap/go-admin/internal/admin/service/Usage"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/pubsub"
	server "github.com/bignyap/go-utilities/server"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminHandler struct {
	BillingService      service.BillingService
	OrganizationService organization.OrganizationService
	PricingService      pricing.PricingService
	ResourceService     resource.ResourceService
	SubscriptionService subscription.SubscriptionService
	UsageService        usage.UsageSummaryService
	DashboardService    dashboard.DashboardService
	ResponseWriter      *server.ResponseWriter
	Logger              api.Logger
	Validator           *validator.Validate
	PubSubClient        pubsub.PubSubClient
}

func NewAdminHandler(
	logger api.Logger,
	rw *server.ResponseWriter,
	db *sqlcgen.Queries,
	conn *pgxpool.Pool,
	validator *validator.Validate,
	pubSubClient pubsub.PubSubClient,
) *AdminHandler {

	return &AdminHandler{

		ResponseWriter: server.NewResponseWriter(logger),
		Logger:         logger,
		Validator:      validator,
		PubSubClient:   pubSubClient,

		BillingService: service.BillingService{
			Logger:       logger,
			Validator:    validator,
			DB:           db,
			Conn:         conn,
			PubSubClient: pubSubClient,
		},
		OrganizationService: organization.OrganizationService{
			Logger:       logger,
			Validator:    validator,
			DB:           db,
			Conn:         conn,
			PubSubClient: pubSubClient,
		},
		PricingService: pricing.PricingService{
			Logger:       logger,
			Validator:    validator,
			DB:           db,
			Conn:         conn,
			PubSubClient: pubSubClient,
		},
		ResourceService: resource.ResourceService{
			Logger:       logger,
			Validator:    validator,
			DB:           db,
			Conn:         conn,
			PubSubClient: pubSubClient,
		},
		SubscriptionService: subscription.SubscriptionService{
			Logger:       logger,
			Validator:    validator,
			DB:           db,
			Conn:         conn,
			PubSubClient: pubSubClient,
		},
		UsageService: usage.UsageSummaryService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
		DashboardService: dashboard.DashboardService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
	}
}
