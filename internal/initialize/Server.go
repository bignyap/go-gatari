package initialize

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/server"
)

func InitializeWebServer(serverType server.ServerType, logger api.Logger, srvc server.Handler) error {

	srvLogger := logger.WithComponent("server.InitializeWebServer")

	srvLogger.Info("Starting")

	config := server.DefaultConfig(serverType)
	ensureDefaultServerConfig(config)

	srvLogger.Info("Configs",
		api.Field{
			Key:   "Port",
			Value: config.Port,
		},
		api.Field{
			Key:   "Environment",
			Value: config.Environment,
		},
		api.Field{
			Key:   "Server Type",
			Value: serverType,
		},
	)

	var srv server.Server
	switch serverType {
	case server.ServerHTTP:
		srv = server.NewHTTPServer(config,
			server.WithLogger(logger),
			server.WithHandler(srvc),
		)
	case server.ServerGRPC:
		srv = server.NewGRPCServer(config,
			server.WithLogger(logger),
			server.WithHandler(srvc),
		)
	default:
		return fmt.Errorf("unsupported server type: %s", serverType)
	}

	if err := srv.Start(); err != nil {
		return fmt.Errorf("error starting the server %s", err)
	}

	srvLogger.Info("Completed")

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
