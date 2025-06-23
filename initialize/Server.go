package initialize

import (
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
	router.RegisterHandlers(
		server.Router().Group("/admin"),
		s.Logger, s.ResponseWriter, s.DB, s.Conn, s.Validator,
	)
	s.Logger.Info("AuthService setup completed")
	return nil
}

func (s *AdminService) Shutdown() error {
	s.Logger.Info("AuthService shutdown initiated")
	return nil
}

func InitializeWebServer(logger api.Logger, conn *pgxpool.Pool, validator *validator.Validate) error {

	config := server.DefaultConfig()
	config.Port = "8081"
	config.Environment = "dev"
	config.Version = "1.0.0"

	adminService := NewAdminService(
		logger, conn, validator,
	)

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
