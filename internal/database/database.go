// database - container for 'dbTX' from 'transaction.go'
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

type Source struct {
	store dbTX
}

func NewSource(db *sql.DB) Source {
	return Source{store: dbTX{DB: db}}
}

// InitDB - create a database connection
//
// 1. get location of file.db
// 2. check file, if not exists -> create database file and install = true
// 3. sql.Open
// 4. if install = true -> create table(s)
func InitDB(cfg *config.Config) (*sql.DB, error) {
	install := false
	if _, err := os.Stat(cfg.DataBaseDataSourceName); err != nil {
		if err := common.CreatePathWithFile(cfg.DataBaseDataSourceName); err != nil {
			return nil, fmt.Errorf("database: file.db create error - %w", err)
		}
		install = true
	}
	db, err := sql.Open("sqlite", cfg.DataBaseDataSourceName)
	if err != nil {
		return nil, fmt.Errorf("database: sql.Open error - %w", err)
	}
	if install {
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeTableCreate)
		defer cancel()
		_, err := db.ExecContext(ctx, schema)
		if err != nil {
			return nil, fmt.Errorf("database: schema init error - %w", err)
		}
	}
	return db, nil
}
