// query - describes requests to database
package database

import "context"

func (s Source) NewTables(ctx context.Context, tables ...string) error {
	newTables := func(ctx context.Context) error {
		for _, table := range tables {
			_, err := s.store.Tx.ExecContext(ctx, table)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return s.store.Transaction(ctx, newTables)
}
