package db

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(
		newDB,
		fx.OnStop(func(db *sqlx.DB) error {
			return db.Close()
		})),
	),
	fx.Invoke(pingDb),
)

func newDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", cfg.DB)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func pingDb(db *sqlx.DB) error {
	return db.Ping()
}
