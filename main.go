package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/transport"
)

func init() {
	if err := godotenv.Load("./init/.env"); err != nil {
		log.Printf("main: no .env file error - %v", err)
	}
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("main: error - %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("main: sql.DB.Close error - %v", err)
		}
	}()
	r := transport.NewTransport(http.NewServeMux())
	if err := r.Run(database.NewSource(db)); err != nil {
		log.Printf("main: error - %v", err)
	}
}
