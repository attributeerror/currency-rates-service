package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	libsql "github.com/libsql/go-libsql"
)

func InitTursoDatabase(options ...TursoDbOption) (*tursoDatabase, error) {
	var _ Database = (*tursoDatabase)(nil)

	opts := &TursoDbOptions{}
	for _, option := range options {
		option(opts)
	}

	if opts.PrimaryUrl == "" {
		return nil, &MissingRequiredOptionError{
			optionName: "PrimaryUrl",
		}
	}

	if opts.AuthToken == "" {
		return nil, &MissingRequiredOptionError{
			optionName: "AuthToken",
		}
	}

	if opts.DbName == "" {
		return nil, &MissingRequiredOptionError{
			optionName: "DbName",
		}
	}

	if opts.SyncInterval.Abs().Milliseconds() == 0 {
		return nil, &MissingRequiredOptionError{
			optionName: "SyncInterval",
		}
	}

	if opts.TableName == "" {
		opts.TableName = "currency_rates"
	}

	dbInst := &tursoDatabase{
		tableName: opts.TableName,
	}

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	dbInst.localDirectory = dir

	connector, err := libsql.NewEmbeddedReplicaConnectorWithAutoSync(dir+"/"+opts.DbName+".db", opts.PrimaryUrl, opts.AuthToken, opts.SyncInterval)
	if err != nil {
		return nil, fmt.Errorf("error whilst connecting to database: %w", err)
	}

	dbInst.connector = connector
	dbInst.db = sql.OpenDB(connector)

	return dbInst, nil
}

func (db *tursoDatabase) GetEuroRateForCurrency(code string) (*float64, error) {
	var euroRate float64
	var euroString string

	sqlStmt := fmt.Sprintf("SELECT printf('%%.5f', to_euro_rate) FROM %s WHERE code = '%s'", db.tableName, code)
	rows, err := db.db.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error whilst querying '%s' table: %w", db.tableName, err)
	}

	defer func() {
		if closeError := rows.Close(); closeError != nil {
			fmt.Println("error closing rows", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	for rows.Next() {
		err = rows.Scan(&euroString)
		if err != nil {
			return nil, err
		}

		euroRate, err = strconv.ParseFloat(euroString, 64)
		if err != nil {
			return nil, fmt.Errorf("error whilst parsing float: %w", err)
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &euroRate, nil
}

func WithPrimaryUrl(primaryUrl string) TursoDbOption {
	return func(opt *TursoDbOptions) {
		opt.PrimaryUrl = primaryUrl
	}
}

func WithAuthToken(authToken string) TursoDbOption {
	return func(opt *TursoDbOptions) {
		opt.AuthToken = authToken
	}
}

func WithDbName(dbName string) TursoDbOption {
	return func(opt *TursoDbOptions) {
		opt.DbName = dbName
	}
}

func WithSyncInterval(interval time.Duration) TursoDbOption {
	return func(opt *TursoDbOptions) {
		opt.SyncInterval = interval
	}
}

func WithTableName(tableName string) TursoDbOption {
	return func(opt *TursoDbOptions) {
		opt.TableName = tableName
	}
}
