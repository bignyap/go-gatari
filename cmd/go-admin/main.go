package main

import (
	"log"
	"os"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/initialize"
	"github.com/bignyap/go-admin/router"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/logger/config"
	"github.com/bignyap/go-utilities/logger/factory"
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
}

func NewAdminService(
	logger api.Logger,
	conn *pgxpool.Pool,
	validator *validator.Validate,
) *AdminService {
	return &AdminService{
		Logger:    logger,
		Validator: validator,
		DB:        sqlcgen.New(conn),
		Conn:      conn,
	}
}

func (s *AdminService) Setup(server server.Server) error {

	setupLogger := s.Logger.WithComponent("server.Setup")

	setupLogger.Info("Starting")

	s.ResponseWriter = server.GetResponseWriter()

	router.RegisterAdminHandlers(
		server.Router(),
		s.Logger, s.ResponseWriter, s.DB, s.Conn, s.Validator,
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

func main() {

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

	validator := validator.New()

	adminSrvc := NewAdminService(
		logger, conn, validator,
	)

	if err := initialize.InitializeWebServer(logger, adminSrvc); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
