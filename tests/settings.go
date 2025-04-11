package tests

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../init/.env"); err != nil {
		log.Printf("settings: no .env file - %v", err)
	}
	if err := os.Setenv("TODO_DBFILE", filepath.Join("../", os.Getenv("TODO_DBFILE"))); err != nil {
		log.Printf("settings: os.Setenv(TODO_DBFILE) - %v", err)
	}
}

var (
	Port         = 7540
	DBFile       = "../scheduler.db"
	FullNextDate = true
	Search       = true
	Token        = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb250ZW50IjoiVGFzayBBY2Nlc3MiLCJleHBsb3JhdGlvbiI6MTc0NDk2NDI2NH0.5KCkfZJwbFoB--pLNGe-Qat8_hVKGtUvMv5TC3zHkn8`
)
