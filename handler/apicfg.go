package handler

import (
	"database/sql"

	"github.com/bignyap/go-admin/database/sqlcgen"
)

type ApiConfig struct {
	DB   *sqlcgen.Queries
	Conn *sql.DB
}

const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusInternalServerError = 500
	StatusNoContent           = 204
)
