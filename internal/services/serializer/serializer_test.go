package serializer

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/config"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/lib/jwtsign"
	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
)

func TestTokenEncode_Response(t *testing.T) {
	cfg, err := config.NewConfig(filepath.Join("..", "..", "..", "init", ".env"))
	require.NoError(t, err, fmt.Sprintf("serializer_test: config error - %v", err))
	require.NoError(t, jwtsign.NewSecretKey(cfg), fmt.Sprintf("serializer_test: secret key error - %v", err))

	serialize := TokenEncode{Content: "Task Access"}
	response, err := serialize.Response()
	require.NoError(t, err, fmt.Sprintf("serializer_test: response error - %v", err))
	assert.Regexp(t, `[a-zA-Z0-9-_.]{148}`, response)
}

func newTask() model.TaskModel {
	return model.TaskModel{
		ID:      123,
		Date:    "20251003",
		Title:   "first",
		Comment: "ololo",
		Repeat:  "d 1",
	}
}

func TestTaskEncode_Response(t *testing.T) {
	serialize := TaskEncode{TaskModel: newTask()}
	response := serialize.Response()
	assert.Equal(t, TaskResponse{
		ID:      "123",
		Date:    "20251003",
		Title:   "first",
		Comment: "ololo",
		Repeat:  "d 1",
	}, *response)
}

func TestTaskListEncode_Response(t *testing.T) {
	serialize := TaskListEncode{Tasks: []model.TaskModel{newTask(), newTask(), newTask()}}
	response := serialize.Response()
	assert.Equal(t, TaslListResponse{
		TasksResp: []TaskResponse{
			{
				ID:      "123",
				Date:    "20251003",
				Title:   "first",
				Comment: "ololo",
				Repeat:  "d 1",
			},
			{
				ID:      "123",
				Date:    "20251003",
				Title:   "first",
				Comment: "ololo",
				Repeat:  "d 1",
			},
			{
				ID:      "123",
				Date:    "20251003",
				Title:   "first",
				Comment: "ololo",
				Repeat:  "d 1",
			},
		},
	}, *response)
}
