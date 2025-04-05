package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("./init/.env"); err != nil {
		log.Printf("no .env file - %v", err)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	if err := http.ListenAndServe(":"+os.Getenv("TODO_PORT"), nil); err != nil {
		log.Fatalf("stop error - %v", err)
	}
}
