package adminHandler

import (
	service "github.com/bignyap/go-admin/internal/admin/service/Billing"
	organization "github.com/bignyap/go-admin/internal/admin/service/Organization"
	pricing "github.com/bignyap/go-admin/internal/admin/service/Pricing"
	resource "github.com/bignyap/go-admin/internal/admin/service/Resource"
	subscription "github.com/bignyap/go-admin/internal/admin/service/Subscription"
	usage "github.com/bignyap/go-admin/internal/admin/service/Usage"
	"github.com/bignyap/go-admin/internal/caching"
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
	}
}

type OrganizationHandler struct {
	OrganizationService organization.OrganizationService
	PricingService      pricing.PricingService
	ResourceService     resource.ResourceService
	SubscriptionService subscription.SubscriptionService
	ResponseWriter      *server.ResponseWriter
	Logger              api.Logger
	Validator           *validator.Validate
	CacheController     *caching.CacheController
}

func NewOrganizationHandler(
	logger api.Logger,
	rw *server.ResponseWriter,
	db *sqlcgen.Queries,
	conn *pgxpool.Pool,
	validator *validator.Validate,
	cacheContoller *caching.CacheController,
) *OrganizationHandler {

	return &OrganizationHandler{

		ResponseWriter:  server.NewResponseWriter(logger),
		Logger:          logger,
		Validator:       validator,
		CacheController: cacheContoller,

		OrganizationService: organization.OrganizationService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
		PricingService: pricing.PricingService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
		ResourceService: resource.ResourceService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
		SubscriptionService: subscription.SubscriptionService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
	}
}

type UsageHandler struct {
	OrganizationService organization.OrganizationService
	PricingService      pricing.PricingService
	ResourceService     resource.ResourceService
	SubscriptionService subscription.SubscriptionService
	ResponseWriter      *server.ResponseWriter
	Logger              api.Logger
	Validator           *validator.Validate
	CacheController     *caching.CacheController
}

func NewUsageHandler(
	logger api.Logger,
	rw *server.ResponseWriter,
	db *sqlcgen.Queries,
	conn *pgxpool.Pool,
	validator *validator.Validate,
	cacheContoller *caching.CacheController,
) *UsageHandler {

	return &UsageHandler{

		ResponseWriter:  server.NewResponseWriter(logger),
		Logger:          logger,
		Validator:       validator,
		CacheController: cacheContoller,

		OrganizationService: organization.OrganizationService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
		PricingService: pricing.PricingService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
		ResourceService: resource.ResourceService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
		SubscriptionService: subscription.SubscriptionService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
		},
	}
}
