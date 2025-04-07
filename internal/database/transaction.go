// transaction
//
// dbTX - wrapper to '*sql.DB' and '*sql.Tx'
//
// inside Transaction set tx to *dbTX.Tx
// inside 'some' query use *Source.tx.Tx
package database

import (
	"context"
	"database/sql"
	"log"
)

type dbTX struct {
	*sql.DB
	*sql.Tx
}

func (base *dbTX) Transaction(ctx context.Context, execute func(ctx context.Context) error) error {
	tx, err := base.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("transaction: Rollback error - %v", err)
		}
	}()
	base.Tx = tx

	if err := execute(ctx); err != nil {
		return err
	}
	return tx.Commit()
}
