package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open(config Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.Url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
