package initializer

import (
	"log"
	"os"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-admin/internal/initialize"
	"github.com/bignyap/go-admin/internal/router"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/logger/config"
	"github.com/bignyap/go-utilities/logger/factory"
	"github.com/bignyap/go-utilities/pubsub"
	"github.com/bignyap/go-utilities/server"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminService struct {
	Logger         api.Logger
	ResponseWriter *server.ResponseWriter
	DB             *sqlcgen.Queries
	Conn           *pgxpool.Pool
	Validator      *validator.Validate
	PubSubClient   pubsub.PubSubClient
}

func NewAdminService(
	logger api.Logger,
	conn *pgxpool.Pool,
	validator *validator.Validate,
	pubSubClient pubsub.PubSubClient,
) *AdminService {
	return &AdminService{
		Logger:       logger,
		Validator:    validator,
		DB:           sqlcgen.New(conn),
		Conn:         conn,
		PubSubClient: pubSubClient,
	}
}

func (s *AdminService) Setup(server server.Server) error {

	setupLogger := s.Logger.WithComponent("server.Setup")

	setupLogger.Info("Starting")

	s.ResponseWriter = server.GetResponseWriter()

	router.RegisterAdminHandlers(
		server.Router(),
		s.Logger,
		s.ResponseWriter,
		s.DB,
		s.Conn,
		s.Validator,
		s.PubSubClient,
	)

	setupLogger.Info("Completed")

	return nil
}

func (s *AdminService) Shutdown() error {

	shtLogger := s.Logger.WithComponent("server.Shutdown")

	shtLogger.Info("Starting")

	if s.Conn != nil {
		s.Conn.Close()
		shtLogger.Info("Database connection pool closed")
	}

	// Add any other cleanup logic here if needed (e.g., flushing logs)

	shtLogger.Info("Completed")

	return nil
}

func InitializeAdminServer() {

	if err := initialize.GetEnvVals(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "dev"
	}

	var logConfig config.LogConfig
	if environment == "prod" {
		logConfig = config.ProductionConfig()
	} else {
		logConfig = config.DevelopmentConfig()
	}
	logger, _ := factory.NewLogger(logConfig)

	conn, err := initialize.LoadDBConn()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	pubSubClient, err := initialize.LoadPubSub()
	if err != nil {
		log.Fatalf("Failed to start the pubsub connection: %v", err)
	}
	defer pubSubClient.Close()

	validator := validator.New()

	adminSrvc := NewAdminService(
		logger, conn, validator, pubSubClient,
	)

	if err := initialize.InitializeWebServer(logger, adminSrvc); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
