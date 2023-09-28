package infra

import (
	"context"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
)

var (
	DBModule = fx.Options(
		fx.Provide(NewDB),
		fx.Invoke(RegisterDB),
	)
  //sqlite datatypes
	schema = `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL,
        bio TEXT,
        image TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    );
  `
)

func NewDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", cfg.DB)

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

      _, err := db.Exec(schema)

      if err != nil {
        return err
      }

			return nil
		},
		OnStop: func(_ context.Context) error {
			return db.Close()
		},
	})
}
