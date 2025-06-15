package initialize

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bignyap/go-admin/database/dbconn"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/handler"
	"github.com/bignyap/go-admin/middlewares"
	"github.com/bignyap/go-admin/router"
	"github.com/joho/godotenv"
)

func GetEnvVals() error {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println(".env file not found, skipping loading environment variables from file")
		return nil
	}

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	return nil
}

func LoadDBConn() (*sql.DB, error) {
	dbConfig := dbconn.DBConfig{
		Driver: "mysql",
		// Driver:   os.Getenv("Driver"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	poolProperties := dbconn.DefaultDBPoolProperties()

	conn, err := dbconn.DBConn(dbConfig, poolProperties)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database %v", err)
	}

	return conn, nil
}

func InitializeWebServer(apiConfig *handler.ApiConfig) error {

	mux := http.NewServeMux()

	router.RegisterHandlers(mux, apiConfig)

	corsMux := middlewares.CorsMiddleware(mux)

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port

	err := http.ListenAndServe(address, corsMux)
	if err != nil {
		return fmt.Errorf("error starting web server: %v", err)
	}
	return nil
}

func InitializeApp() {

	if err := GetEnvVals(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	conn, err := LoadDBConn()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	apiCfg := &handler.ApiConfig{
		DB:   sqlcgen.New(conn),
		Conn: conn,
	}

	if err := InitializeWebServer(apiCfg); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
