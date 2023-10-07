package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SqlStore struct {
	Db *sqlx.DB
}

func (s *SqlStore) CreateTx(fn func(*sqlx.Tx) error) error {
	tx, err := s.Db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
