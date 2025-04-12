// route - implementation of http.HandlerFunc
package transport

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/deserializer"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services/serializer"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// ErrTransportInvalidParam - invalid param from r.URL.Query
var ErrTransportInvalidParam = errors.New("invalid param")

func Login(w http.ResponseWriter, r *http.Request) {
	deserialize := deserializer.NewLoginDecode()
	if err := deserialize.Decode(r); err != nil {
		common.EncodeJSON(w, http.StatusUnprocessableEntity, common.Message{"error": err.Error()})
		return
	}
	password := os.Getenv("TODO_PASSWORD")
	login := deserialize.Model()
	if !login.ValidPassword(password) {
		common.EncodeJSON(w, http.StatusForbidden, common.Message{"error": model.ErrModelsLoginInvalidPassword.Error()})
		return
	}
	serialize := serializer.TokenEncode{Content: "Task Access"}
	tokenResponse, err := serialize.Response()
	if err != nil {
		common.EncodeJSON(w, http.StatusInternalServerError, common.Message{"error": err.Error()})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenResponse.Token,
		Expires: time.Now().UTC().Add(7 * 24 * time.Hour),
	})
	common.EncodeJSON(w, http.StatusOK, tokenResponse)
}

func TaskNew(db model.TaskCreate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deserialize := deserializer.NewTaskDecode()
		if err := deserialize.Decode(r, services.NextDate); err != nil {
			common.EncodeJSON(w, http.StatusUnprocessableEntity, common.Message{"error": err.Error()})
			return
		}
		taskID, err := db.SaveOneTask(r.Context(), deserialize.Model())
		if err != nil {
			common.EncodeJSON(w, http.StatusConflict, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(w, http.StatusCreated, common.Message{"id": strconv.Itoa(int(taskID))})
	}
}

func TaskRetrive(db model.TaskRead) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.Message{"error": ErrTransportInvalidParam.Error()})
			return
		}
		task, err := db.FindOneTask(r.Context(), uint(id))
		if err != nil {
			common.EncodeJSON(w, http.StatusNotFound, common.Message{"error": err.Error()})
			return
		}
		serialize := serializer.TaskEncode{TaskModel: task}
		common.EncodeJSON(w, http.StatusOK, serialize.Response())
		//common.EncodeJSON( w, http.StatusOK, common.Message{"task": serialize.Response()})
	}
}

func TaskChange(db model.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deserialize := deserializer.NewTaskDecode()
		if err := deserialize.Decode(r, services.NextDate); err != nil {
			common.EncodeJSON(w, http.StatusUnprocessableEntity, common.Message{"error": err.Error()})
			return
		}
		task := deserialize.Model()
		if task.ID == 0 {
			common.EncodeJSON(w, http.StatusUnprocessableEntity, common.Message{"error": deserializer.ErrServicesWrongID.Error()})
			return
		}
		if err := db.NewDataTask(r.Context(), task); err != nil {
			common.EncodeJSON(w, http.StatusNotFound, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{})
	}
}

func TaskRemove(db model.TaskDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.Message{"error": ErrTransportInvalidParam.Error()})
			return
		}
		if err := db.ExpirationTask(r.Context(), uint(id)); err != nil {
			common.EncodeJSON(w, http.StatusNotFound, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{})
	}
}

func TaskDone(db multiTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			common.EncodeJSON(w, http.StatusBadRequest, common.Message{"error": ErrTransportInvalidParam.Error()})
			return
		}
		ctx := r.Context()
		task, err := db.FindOneTask(ctx, uint(id))
		if err != nil {
			common.EncodeJSON(w, http.StatusNotFound, common.Message{"error": err.Error()})
			return
		}
		if err := task.UpdateDate(services.NextDate); err != nil {
			if errors.Is(err, model.ErrModelTaskDone) {
				if err := db.ExpirationTask(ctx, task.ID); err != nil {
					common.EncodeJSON(w, http.StatusInternalServerError, common.Message{"error": err.Error()})
					return
				}
				common.EncodeJSON(w, http.StatusOK, common.Message{})
				return
			}
			common.EncodeJSON(w, http.StatusInternalServerError, common.Message{"error": err.Error()})
			return
		}
		if err := db.NewDataTask(ctx, task); err != nil {
			common.EncodeJSON(w, http.StatusInternalServerError, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(w, http.StatusOK, common.Message{})
	}
}

func TaskRetriveList(db model.TaskRead) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("search")
		taskProperty := services.NewTaskProperty(param, 123)
		tasks, err := db.FindTaskList(r.Context(), taskProperty)
		if err != nil {
			common.EncodeJSON(w, http.StatusInternalServerError, common.Message{"error": err.Error()})
			return
		}
		serialize := serializer.TaskListEncode{Tasks: tasks}
		common.EncodeJSON(w, http.StatusOK, common.Message{"tasks": serialize.Response()})
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
	newDate, err := services.NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(newDate))
	if err != nil {
		log.Printf("route: http.ResponseWriter.Write error - %v", err)
	}
}
