package database

import (
	"database/sql"
	"time"

	"github.com/libsql/go-libsql"
)

type (
	Database interface {
		GetEuroRateForCurrency(code string) (*float64, error)
	}

	TursoDatabase interface {
		Database
	}

	tursoDatabase struct {
		db             *sql.DB
		connector      *libsql.Connector
		tableName      string
		localDirectory string
	}

	TursoDbOptions struct {
		PrimaryUrl   string
		AuthToken    string
		DbName       string
		TableName    string
		SyncInterval time.Duration
	}

	TursoDbOption func(*TursoDbOptions)
)
