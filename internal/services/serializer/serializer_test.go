package serializer

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

func init() {
	if err := godotenv.Load("../../../init/.env"); err != nil {
		log.Printf("autorization_test: no .env file - %v", err)
	}
	common.SecretKey = os.Getenv("TODO_SECRET_KEY")
	if common.SecretKey == "" {
		log.Printf("autorization_test: SecretKey is empty")
	}
}

func TestTokenEncode_Response(t *testing.T) {
	serialize := TokenEncode{Content: "Task Access"}
	response, err := serialize.Response()
	require.NoError(t, err, fmt.Sprintf("Response error - %v", err))
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
	}, response)
}

func TestTaskListEncode_Response(t *testing.T) {
	serialize := TaskListEncode{Tasks: []model.TaskModel{newTask(), newTask(), newTask()}}
	response := serialize.Response()
	assert.Equal(t, []TaskResponse{
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
	}, response)
}
