// app - heart of application
//
// container for all services see (/internal/services.services.go)
package app

import (
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/usecase"
)

// main object of application
type Sheduler struct {
	usecase.TaskService

	usecase.LoginService

	usecase.AuthService
}

func NewSheduler(
	cfg *config.Config,
	taskStore usecase.MultiTask,
	loginStore usecase.MultiLogin) (Sheduler, error) {
	taskService, err := usecase.NewTaskService(cfg, taskStore)
	if err != nil {
		return Sheduler{}, err
	}
	loginService := usecase.NewLoginService(loginStore)
	authService := usecase.NewAuthService()
	return Sheduler{
		TaskService:  taskService,
		LoginService: loginService,
		AuthService:  authService,
	}, nil
}
