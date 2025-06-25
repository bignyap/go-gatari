package resource

import (
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResourceService struct {
	DB        *sqlcgen.Queries
	Conn      *pgxpool.Pool
	Logger    api.Logger
	Validator *validator.Validate
}
