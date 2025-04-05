package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/server"
)

func init() {
	if err := godotenv.Load("./init/.env"); err != nil {
		log.Printf("main: no .env file error - %v", err)
	}
}

func main() {
	db := database.InitDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("main: sql.DB.Close error - %v", err)
		}
	}()
	_ = database.NewSource(db)
	mux := http.NewServeMux()
	srv := server.InitSRV(mux)

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	srv.ListenAndServeAndShut(server.ServerTimeoutShut)
}
