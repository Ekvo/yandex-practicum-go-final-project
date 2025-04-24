// run - describe process initializing the application and starting server
package app

import (
	"log"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/database"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/datauser"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/jwtsign"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/transport"
)

// 1. set secretkey for jwt.Token  -> 'jwtsign.NewSecretKey'
// 2. open database                -> 'database.InitDB'
// 3. create Sheduler heart of app -> 'NewSheduler'
// 4. create server and router     -> `transport.NewTransport`
// 5. start (close inside)         -> `Start`
func Run(cfg *config.Config) {
	if err := jwtsign.NewSecretKey(cfg); err != nil {
		log.Fatalf("app: error - %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("app: error - %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("app: sql.DB.Close error - %v", err)
		}
	}()

	sheduler, err := NewSheduler(
		cfg,
		database.NewSource(db),
		datauser.NewUserData(cfg))
	if err != nil {
		log.Fatalf("app: error - %v", err)
	}

	r := transport.NewTransport(cfg)

	if err := r.Start(sheduler); err != nil {
		log.Printf("app: error - %v", err)
	}
}
