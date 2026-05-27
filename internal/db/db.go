package db

import (
	"errors"
	"fmt"
	"server/internal/config"

	"github.com/jackc/pgx"
)

var ErrDBConnection = errors.New("database: connection failed")

// NewDB returns a new connection or connection error
func NewDB(cfg *config.Config) (*pgx.Conn, error) {

	sslMode := "disable"
	if cfg.Env == "production" {
	  sslMode = "require"
	}
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
  cfg.DbUser, cfg.DbPwd, cfg.DbHost, cfg.DbPort, cfg.DbName, sslMode)
	
	// Parses the URI and returns ConnConfig type or an error
	dbConfig, parsingErr := pgx.ParseURI(connStr)
	if parsingErr != nil {
		return nil, fmt.Errorf("invalid connection string: %w ", parsingErr)
	}

	// Establishes a psql-db connection or returns an error
	conn, connErr := pgx.Connect(dbConfig)
	if connErr != nil {
		return nil, fmt.Errorf("%w: %v ", ErrDBConnection, connErr)
	} else {
		return conn, nil
	}
}
