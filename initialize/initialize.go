package initialize

import (
	"log"

	"github.com/bignyap/go-utilities/logger/api"
	"github.com/go-playground/validator"
)

func InitializeApp() {

	logger := &api.DefaultLogger{}

	if err := GetEnvVals(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

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
