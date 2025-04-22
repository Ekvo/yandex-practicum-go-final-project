// route - bizcase of http.HandlerFunc
package transport

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/nextdate"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/deserializer"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/entity"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/usecase"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// ErrTransportInvalidParam - invalid param from r.URL.Query
var ErrTransportInvalidParam = errors.New("invalid param")

func Login(loginService services.LoginValidPasswordCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deserialize := deserializer.NewLoginDecode()
		if err := deserialize.Decode(r); err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.NewError(err))
			return
		}
		token, err := loginService.CreateToken(r.Context(), deserialize.Model())
		if err != nil {
			code := 0
			if errors.Is(err, usecase.ErrCaseLoginNotFound) {
				code = http.StatusForbidden
			} else if errors.Is(err, services.ErrServicesInternalError) {
				code = http.StatusInternalServerError
			} else {
				code = http.StatusUnprocessableEntity
			}
			common.EncodeJSON(w, code, common.NewError(err))
			return
		}
		//http.SetCookie(w, &http.Cookie{ //for  postman
		//	Name:    `token`,
		//	Value:   token,
		//	Expires: time.Now().Add(time.Hour * 24 * 7),
		//})
		common.EncodeJSON(w, http.StatusOK, token)
	}
}

func TaskNew(taskService services.TaskCreateCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deserialize := deserializer.NewTaskDecode()
		if err := deserialize.Decode(r); err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.NewError(err))
			return
		}
		taskID, err := taskService.CreateTask(r.Context(), deserialize.Model())
		if err != nil {
			code := 0
			if errors.Is(err, usecase.ErrCaseTaskAlreadyExist) {
				code = http.StatusConflict
			} else if errors.Is(err, services.ErrServicesInternalError) {
				code = http.StatusInternalServerError
			} else {
				code = http.StatusUnprocessableEntity
			}
			common.EncodeJSON(w, code, common.NewError(err))
			return
		}
		common.EncodeJSON(w, http.StatusCreated, taskID)
	}
}

func TaskRetrieve(taskService services.TaskReadCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.NewError(ErrTransportInvalidParam))
			return
		}
		task, err := taskService.ReadTask(r.Context(), uint(id))
		if err != nil {
			code := 0
			if errors.Is(err, usecase.ErrCaseTaskNotFound) {
				code = http.StatusNotFound
			} else if errors.Is(err, services.ErrServicesInternalError) {
				code = http.StatusInternalServerError
			} else {
				code = http.StatusUnprocessableEntity
			}
			common.EncodeJSON(w, code, common.NewError(err))
			return
		}
		common.EncodeJSON(w, http.StatusOK, task)
	}
}

func TaskChange(taskService services.TaskUpdateCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deserialize := deserializer.NewTaskDecode()
		if err := deserialize.Decode(r); err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.NewError(err))
			return
		}
		if err := taskService.UpdateTask(r.Context(), deserialize.Model()); err != nil {
			code := 0
			if errors.Is(err, usecase.ErrCaseTaskNotFound) {
				code = http.StatusNotFound
			} else if errors.Is(err, services.ErrServicesInternalError) {
				code = http.StatusInternalServerError
			} else {
				code = http.StatusUnprocessableEntity
			}
			common.EncodeJSON(w, code, common.NewError(err))
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{})
	}
}

func TaskRemove(taskService services.TaskDeleteCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.NewError(ErrTransportInvalidParam))
			return
		}
		if err := taskService.DeleteTask(r.Context(), uint(id)); err != nil {
			code := 0
			if errors.Is(err, usecase.ErrCaseTaskNotFound) {
				code = http.StatusNotFound
			} else if errors.Is(err, services.ErrServicesInternalError) {
				code = http.StatusInternalServerError
			} else {
				code = http.StatusUnprocessableEntity
			}
			common.EncodeJSON(w, code, common.NewError(err))
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{})
	}
}

func TaskDone(taskService services.TaskDoneCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.NewError(ErrTransportInvalidParam))
			return
		}
		if err := taskService.DoneTask(r.Context(), uint(id)); err != nil {
			code := 0
			if errors.Is(err, usecase.ErrCaseTaskNotFound) {
				code = http.StatusNotFound
			} else if errors.Is(err, services.ErrServicesInternalError) {
				code = http.StatusInternalServerError
			} else {
				code = http.StatusUnprocessableEntity
			}
			common.EncodeJSON(w, code, common.NewError(err))
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{})
	}
}

func TaskRetriveList(taskService services.TaskReadCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("search")
		taskProperty := entity.NewTaskProperty(param, 123)
		tasks, err := taskService.ReadTaskList(r.Context(), taskProperty)
		if err != nil {
			common.EncodeJSON(w, http.StatusInternalServerError, common.NewError(err))
			return
		}
		common.EncodeJSON(w, http.StatusOK, tasks)
	}
}

func TestNextDate(w http.ResponseWriter, r *http.Request) {
	timeNowStr := r.URL.Query().Get("now")
	dstart := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	now, err := time.Parse(model.DateFormat, timeNowStr)
	if err != nil && timeNowStr != "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if timeNowStr == "" {
		now = time.Now()
	}
	newDate, err := nextdate.NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(newDate))
	if err != nil {
		log.Printf("route: http.ResponseWriter.Write error - %v", err)
	}
}
