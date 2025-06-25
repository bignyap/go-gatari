package initialize

import (
	"fmt"
	"os"

	"github.com/bignyap/go-admin/database/dbconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bignyap/go-utilities/database"
)

func LoadDBConn() (*pgxpool.Pool, error) {

	dbConfig := &database.ConnectionString{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Options: map[string]string{
			"sslmode": "disable",
		},
	}

	poolProperties := database.DefaultPoolConfig()

	dbConn, err := dbconn.DBConn(
		"go-admin", "postgres",
		dbConfig, poolProperties,
	)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database %v", err)
	}

	return dbConn.Connection.GetPgxPool(), nil
}
