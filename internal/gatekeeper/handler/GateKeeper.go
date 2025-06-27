package gateKeeperHandler

import (
	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/logger/api"
	server "github.com/bignyap/go-utilities/server"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GateKeeperHandler struct {
	// OrganizationService organization.OrganizationService
	// PricingService      pricing.PricingService
	// ResourceService     resource.ResourceService
	// SubscriptionService subscription.SubscriptionService
	ResponseWriter  *server.ResponseWriter
	Logger          api.Logger
	Validator       *validator.Validate
	CacheController *caching.CacheController
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

		// OrganizationService: organization.OrganizationService{
		// 	Logger:    logger,
		// 	Validator: validator,
		// 	DB:        db,
		// 	Conn:      conn,
		// },
		// PricingService: pricing.PricingService{
		// 	Logger:    logger,
		// 	Validator: validator,
		// 	DB:        db,
		// 	Conn:      conn,
		// },
		// ResourceService: resource.ResourceService{
		// 	Logger:    logger,
		// 	Validator: validator,
		// 	DB:        db,
		// 	Conn:      conn,
		// },
		// SubscriptionService: subscription.SubscriptionService{
		// 	Logger:    logger,
		// 	Validator: validator,
		// 	DB:        db,
		// 	Conn:      conn,
		// },
	}
}
