package subscription

import (
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/pubsub"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionService struct {
	DB           *sqlcgen.Queries
	Conn         *pgxpool.Pool
	Logger       api.Logger
	Validator    *validator.Validate
	PubSubClient pubsub.PubSubClient
}
