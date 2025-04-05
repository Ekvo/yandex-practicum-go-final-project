// database - container for 'dbTX' from 'transaction.go'
package database

import (
	"context"
	"database/sql"
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

func InitDB() *sql.DB {
	pathDB := os.Getenv("TODO_DBFILE")
	install := false
	if _, err := os.Stat(pathDB); err != nil {
		if err := common.CreatePathWithFile(pathDB); err != nil {
			log.Fatalf("database: file.db create error - %v", err)
		}
		install = true
	}
	db, err := sql.Open("sqlite", pathDB)
	if err != nil {
		log.Fatalf("database: sql.Open error - %v", err)
	}
	if install {
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeTableCreate)
		defer cancel()
		_, err := db.ExecContext(ctx, Schema)
		if err != nil {
			log.Fatalf("database: schema init error - %v", err)
		}
	}
	return db
}
