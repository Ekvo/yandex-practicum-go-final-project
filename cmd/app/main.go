package main

import (
	"log"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/app"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
)

func main() {
	cfg, err := config.NewConfig("./init/.env")
	if err != nil {
		log.Fatalf("main: error - %v", err)
	}
	app.Run(cfg)
}
