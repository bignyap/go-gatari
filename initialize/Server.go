package initialize

import (
	"os"
	"strconv"
	"time"

	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/router"
	"github.com/bignyap/go-utilities/logger/api"
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

	s.ResponseWriter = server.GetResponseWriter()

	adminGrp := server.Router().Group("/admin")

	router.RegisterHandlers(
		adminGrp,
		s.Logger, s.ResponseWriter, s.DB, s.Conn, s.Validator,
	)
	s.Logger.Info("AuthService setup completed")
	return nil
}

func (s *AdminService) Shutdown() error {
	s.Logger.Info("AuthService shutdown initiated")

	if s.Conn != nil {
		s.Conn.Close()
		s.Logger.Info("Database connection pool closed")
	}

	// Add any other cleanup logic here if needed (e.g., flushing logs)

	s.Logger.Info("AuthService shutdown completed")
	return nil
}

func InitializeWebServer(logger api.Logger, conn *pgxpool.Pool, validator *validator.Validate) error {

	adminService := NewAdminService(
		logger, conn, validator,
	)

	config := server.DefaultConfig()
	ensureDefaultServerConfig(config)

	s := server.NewHTTPServer(
		config,
		server.WithLogger(logger),
		server.WithHandler(adminService),
	)

	if err := s.Start(); err != nil {
		logger.Error("Server failed", err)
	}

	return nil
}

func ensureDefaultServerConfig(config *server.Config) {

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8080"
	}
	config.Port = port

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "dev"
	}
	config.Environment = environment

	version := os.Getenv("VERSION")
	if version == "" {
		version = "UNDEFINED"
	}
	config.Version = version

	maxRequestSize := os.Getenv("MAX_REQUEST_SIZE")
	if maxRequestSize == "" {
		config.MaxRequestSize = 10 * 1024 * 1024 // Default to 10 MB
	} else {
		size, err := strconv.ParseInt(maxRequestSize, 10, 64)
		if err != nil {
			config.MaxRequestSize = 10 * 1024 * 1024 // Default to 10 MB
		} else {
			config.MaxRequestSize = size * 1024 * 1024
		}
	}

	enableProfiling := os.Getenv("ENABLE_PROFILING")
	if enableProfiling == "" {
		config.EnableProfiling = false
	} else {
		profiling, err := strconv.ParseBool(enableProfiling)
		if err != nil {
			config.EnableProfiling = false
		} else {
			config.EnableProfiling = profiling
		}
	}

	shutdownTimeout := os.Getenv("SHUTDOWN_TIMEOUT")
	if shutdownTimeout == "" {
		config.ShutdownTimeout = 30 * time.Second // Default to 30 seconds
	} else {
		timeout, err := time.ParseDuration(shutdownTimeout)
		if err != nil {
			config.ShutdownTimeout = 30 * time.Second // Default to 30 seconds
		} else {
			config.ShutdownTimeout = timeout
		}
	}
}
