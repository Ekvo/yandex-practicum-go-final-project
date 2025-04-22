// services - describes all biz logic of application
//
// contaion only interface and variable of Error

package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/entity"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/serializer"
)

// ErrServicesInternalError - change internals erros of application on ErrServicesInternalError
var ErrServicesInternalError = errors.New("internal error")

type (
	// TaskCreateCase - logic of create Task
	TaskCreateCase interface {
		CreateTask(
			ctx context.Context,
			task model.TaskModel) (*serializer.TaskIDResponse, error)
	}

	// TaskReadCase  - logic of read task(s)
	TaskReadCase interface {
		ReadTask(
			ctx context.Context,
			id uint) (*serializer.TaskResponse, error)
		ReadTaskList(
			ctx context.Context,
			property *entity.TaskProperty) (*serializer.TaslListResponse, error)
	}

	// TaskUpdateCase - logic of updatre Task
	TaskUpdateCase interface {
		UpdateTask(ctx context.Context, task model.TaskModel) error
	}

	// TaskDeleteCase - biz logic of delete Task
	TaskDeleteCase interface {
		DeleteTask(ctx context.Context, id uint) error
	}

	// TaskDoneCase - logic for task marked Done
	TaskDoneCase interface {
		DoneTask(ctx context.Context, id uint) error
	}

	// LoginValidPasswordCase - logic of login fro application
	LoginValidPasswordCase interface {
		CreateToken(
			ctx context.Context,
			login model.LoginModel) (*serializer.TokenResponse, error)
		UserExist(ctx context.Context) (bool, error)
	}

	// AutorizationCase - logic of autorization
	AutorizationCase interface {
		AuthZ(r *http.Request) error
	}
)
