package initialize

import (
	"log"
	"os"

	"github.com/bignyap/go-utilities/logger/config"
	"github.com/bignyap/go-utilities/logger/factory"
	"github.com/go-playground/validator"
)

func InitializeApp() {

	if err := GetEnvVals(); err != nil {
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

	conn, err := LoadDBConn()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	validator := validator.New()

	if err := InitializeWebServer(logger, conn, validator); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
