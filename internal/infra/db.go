package infra

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"

	"github.com/berkeleytrue/conduit/config"
)

func NewDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", cfg.DB)

	if err != nil {
		return nil, err
	}

	return db, nil
}
