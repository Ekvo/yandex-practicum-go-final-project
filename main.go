package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/server"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
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

	mux.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		now := r.URL.Query().Get("now")
		dstart := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")
		t := time.Time{}
		if now == "" {
			t = time.Now()
		} else {
			m, err := time.Parse("20060102", now)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			t = m
		}
		newt, err := services.NextDate(t, dstart, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(newt))
	})

	srv.ListenAndServeAndShut(server.ServerTimeoutShut)
}
