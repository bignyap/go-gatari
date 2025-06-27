package dbconn

import (
	"fmt"

	"github.com/bignyap/go-utilities/database"
)

func DBConn(
	name string,
	driverStr string,
	cs *database.ConnectionString,
	pool *database.ConnectionPoolConfig,
) (*database.Database, error) {

	driver, err := database.ParseDriver(driverStr)
	if err != nil {
		return nil, fmt.Errorf("invalid driver: %w", err)
	}

	if pool == nil {
		pool = database.DefaultPoolConfig()
	}

	db, err := database.NewDatabase(&database.DatabaseConfig{
		Name:                 name,
		Driver:               driver,
		ConnectionString:     cs,
		ConnectionPoolConfig: pool,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	if err := db.Connection.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
