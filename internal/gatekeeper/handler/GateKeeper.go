package gateKeeperHandler

import (
	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	"github.com/bignyap/go-utilities/logger/api"
	server "github.com/bignyap/go-utilities/server"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GateKeeperHandler struct {
	GateKeepingService *gatekeeping.GateKeepingService
	ResponseWriter     *server.ResponseWriter
	Logger             api.Logger
	Validator          *validator.Validate
	Cache              *caching.CacheController
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

		ResponseWriter: server.NewResponseWriter(logger),
		Logger:         logger,
		Validator:      validator,
		Cache:          cacheContoller,

		GateKeepingService: &gatekeeping.GateKeepingService{
			Logger:    logger,
			Validator: validator,
			DB:        db,
			Conn:      conn,
			Cache:     cacheContoller,
		},
	}
}
