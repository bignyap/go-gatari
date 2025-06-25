package usage

import (
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/go-playground/validator"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bignyap/go-utilities/logger/api"
)

type UsageSummaryService struct {
	DB        *sqlcgen.Queries
	Conn      *pgxpool.Pool
	Logger    api.Logger
	Validator *validator.Validate
}
