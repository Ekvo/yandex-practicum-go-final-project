// database - container for 'dbTX' from 'transaction.go'
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"

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
func InitDB() (*sql.DB, error) {
	pathDB := os.Getenv("TODO_DBFILE")
	if pathDB == "" {
		pathDB = "./storage/scheduler.db"
	}
	install := false
	if _, err := os.Stat(pathDB); err != nil {
		if err := common.CreatePathWithFile(pathDB); err != nil {
			log.Fatalf("database: file.db create error - %v", err)
		}
		install = true
	}
	db, err := sql.Open("sqlite", pathDB)
	if err != nil {
		return nil, fmt.Errorf("database: sql.Open error - %v", err)
	}
	if install {
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeTableCreate)
		defer cancel()
		_, err := db.ExecContext(ctx, Schema)
		if err != nil {
			return nil, fmt.Errorf("database: schema init error - %v", err)
		}
	}
	return db, nil
}
