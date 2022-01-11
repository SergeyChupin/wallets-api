package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open(logger *log.Logger, config Config) (*sql.DB, error) {
	logger.Println("Opening postgres - pgx data source")
	db, err := sql.Open("pgx", config.Url)
	if err != nil {
		return nil, fmt.Errorf("postgres - Open - sql.Open: %w", err)
	}
	logger.Println("Trying to connect to postgres database")
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres - Open - db.Ping: %w", err)
	}
	logger.Println("Connect to postgres database successful")
	return db, nil
}
