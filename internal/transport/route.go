// route - implementation of http.HandlerFunc
package transport

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

var ErrTransportInvalidParam = errors.New("invalid param")

func TaskNew(db model.TaskCreate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deserializer := services.NewTaskDecode()
		if err := deserializer.Decode(r, services.NextDate); err != nil {
			common.EncodeJSON(w, http.StatusUnprocessableEntity, common.Message{"error": err.Error()})
			return
		}
		taskID, err := db.SaveOneTask(r.Context(), deserializer.Model())
		if err != nil {
			common.EncodeJSON(w, http.StatusConflict, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(w, http.StatusCreated, common.Message{"id": taskID})
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
		serialize := services.TaskEncode{TaskModel: task}
		common.EncodeJSON(w, http.StatusOK, serialize.Response())
		//common.EncodeJSON( w, http.StatusOK, common.Message{"task": serialize.Response()})
	}
}

func TaskChange(db model.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deserializer := services.NewTaskDecode()
		if err := deserializer.Decode(r, services.NextDate); err != nil {
			common.EncodeJSON(w, http.StatusUnprocessableEntity, common.Message{"error": err.Error()})
			return
		}
		if err := db.NewDataTask(r.Context(), deserializer.Model()); err != nil {
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
		serialize := services.TaskListEncode{Tasks: tasks}
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
		now = time.Now().UTC()
	}
	newDate, err := services.NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(newDate))
}
