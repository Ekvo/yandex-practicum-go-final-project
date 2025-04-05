package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
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

	http.Handle("/", http.FileServer(http.Dir("./web")))
	if err := http.ListenAndServe(":"+os.Getenv("TODO_PORT"), nil); err != nil {
		log.Fatalf("stop error - %v", err)
	}
}
