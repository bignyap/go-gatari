package cachemanagement

import (
	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CacheManagementService struct {
	DB        *sqlcgen.Queries
	Conn      *pgxpool.Pool
	Logger    api.Logger
	Validator *validator.Validate
	Cache     *caching.CacheController
}
