package db

import (
	"errors"
	"fmt"
	"server/internal/config"

	"github.com/jackc/pgx"
)

var ErrDBConnection = errors.New("database: connection failed")

// NewDB returns a new connection
func NewDB(cfg *config.Config) (*pgx.Conn, error) {	
	
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPwd, cfg.DbHost, cfg.DbPort, cfg.DbName)

	// initialize db connection
	config, parsingErr := pgx.ParseURI(connStr)
	if parsingErr != nil {
		return nil, fmt.Errorf("invalid connection string: %w", parsingErr)
	}

	conn, connErr := pgx.Connect(config)
	if connErr != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBConnection, connErr)
	} else {
		return conn, nil
	}
}