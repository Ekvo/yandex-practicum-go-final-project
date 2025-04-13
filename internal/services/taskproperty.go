// taskproperty - describes rules for find 'model.TaskModel' array from database
package services

import (
	"time"

	"golang.org/x/exp/rand"
)

// dateFormatFromParam - valid format date from param (/api/tasks?search=02.01.2006)
const dateFormatFromParam = "02.01.2006"

// range of limit field ”TaskProperty.limit"
const (
	minLimit = 10
	maxLimit = 50
)

type TaskProperty struct {
	// find in base by world : SELECT * FROM table WHERE colum LIKE world;
	word string

	// find by date
	date time.Time

	// use 'LIMIT' when searching in database
	limit uint
}

func NewTaskProperty(property string, limit uint) *TaskProperty {
	taskProperty := &TaskProperty{}
	taskProperty.parseProperty(property)
	taskProperty.setLimit(limit)
	return taskProperty
}

// setLimit - check limit and if not in [10,50] -> create random limit
func (t *TaskProperty) setLimit(limit uint) {
	if limit < minLimit || limit > maxLimit {
		rand.Seed(uint64(time.Now().UTC().UnixNano()))
		limit = uint(rand.Intn(maxLimit-minLimit) + minLimit) // (40)+10 -> [10,50]
	}
	t.limit = limit
}

// parseProperty - if property not empty - create 'date' or 'word'
func (t *TaskProperty) parseProperty(property string) {
	if property == "" {
		return
	}
	date, err := time.Parse(dateFormatFromParam, property)
	if err != nil {
		t.word = property
		return
	}
	t.date = date
}

func (t *TaskProperty) IsDate() bool {
	return !t.date.IsZero()
}

func (t *TaskProperty) IsWord() bool {
	return t.word != ""
}

func (t *TaskProperty) PassDate() time.Time {
	return t.date
}

func (t *TaskProperty) PassWord() string {
	return t.word
}

func (t *TaskProperty) PassLimit() uint {
	return t.limit
}
