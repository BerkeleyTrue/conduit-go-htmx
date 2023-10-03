package db

import (
	"context"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
)

func NewDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", cfg.DB)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func RegisterDB(lc fx.Lifecycle, db *sqlx.DB) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := db.Ping(); err != nil {
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			return db.Close()
		},
	})
}
