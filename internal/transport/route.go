// route - implementation of http.HandlerFunc
package transport

import (
	"net/http"
	"time"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/services"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

func TaskNew(db model.TaskCreate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		taskDecode := services.NewTaskDecode()
		if err := taskDecode.Decode(r, services.NextDate); err != nil {
			common.EncodeJSON(ctx, w, http.StatusUnprocessableEntity, common.Message{"error": err.Error()})
			return
		}
		taskID, err := db.SaveOneTask(ctx, taskDecode.Model())
		if err != nil {
			common.EncodeJSON(ctx, w, http.StatusConflict, common.Message{"error": err.Error()})
			return
		}
		common.EncodeJSON(ctx, w, http.StatusCreated, common.Message{"id": taskID})
	}
}

func TaskList(db model.TaskFind) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		param := r.URL.Query().Get("search")
		taskProperty := services.NewTaskProperty(param, 123)
		tasks, err := db.FindTaskList(ctx, taskProperty)
		if err != nil {
			common.EncodeJSON(ctx, w, http.StatusInternalServerError, common.Message{"error": err.Error()})
			return
		}
		taskListEncode := services.TaskListEncode{Tasks: tasks}
		common.EncodeJSON(ctx, w, http.StatusOK, common.Message{"tasks": taskListEncode.Response()})
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
