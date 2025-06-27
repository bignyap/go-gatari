package gateKeeperHandler

import (
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/internal/caching"
	organization "github.com/bignyap/go-admin/internal/service/Organization"
	pricing "github.com/bignyap/go-admin/internal/service/Pricing"
	resource "github.com/bignyap/go-admin/internal/service/Resource"
	subscription "github.com/bignyap/go-admin/internal/service/Subscription"
	"github.com/bignyap/go-utilities/logger/api"
	server "github.com/bignyap/go-utilities/server"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GateKeeperHandler struct {
	OrganizationService organization.OrganizationService
	PricingService      pricing.PricingService
	ResourceService     resource.ResourceService
	SubscriptionService subscription.SubscriptionService
	ResponseWriter      *server.ResponseWriter
	Logger              api.Logger
	Validator           *validator.Validate
	CacheController     *caching.CacheController
}

func NewGateKeeperHandler(
	logger api.Logger,
	rw *server.ResponseWriter,
	db *sqlcgen.Queries,
	conn *pgxpool.Pool,
	validator *validator.Validate,
	cacheContoller *caching.CacheController,
) *GateKeeperHandler {

	return &GateKeeperHandler{

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
