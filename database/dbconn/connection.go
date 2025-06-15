package dbconn

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type DBPoolProperties struct {
	SetMaxOpenConns    int
	SetMaxIdleConns    int
	SetConnMaxIdleTime int
	SetConnMaxLifetime int
}

func DefaultDBPoolProperties() DBPoolProperties {
	return DBPoolProperties{
		SetMaxOpenConns:    30,
		SetMaxIdleConns:    10,
		SetConnMaxIdleTime: 300,
		SetConnMaxLifetime: 600,
	}
}

func DBConn(
	config DBConfig,
	poolProps DBPoolProperties,
) (*sql.DB, error) {

	driverName := config.Driver
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&allowAllFiles=true",
		config.User, config.Password,
		config.Host, config.Port,
		config.DBName,
	)

	connPool, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}

	err = connPool.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	connPool.SetMaxOpenConns(poolProps.SetMaxOpenConns)
	connPool.SetMaxIdleConns(poolProps.SetMaxIdleConns)
	connPool.SetConnMaxIdleTime(
		time.Duration(poolProps.SetConnMaxIdleTime) * time.Second,
	)
	connPool.SetConnMaxLifetime(
		time.Duration(poolProps.SetConnMaxLifetime) * time.Second,
	)

	return connPool, nil
}
