package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskProperty(t *testing.T) {
	property := NewTaskProperty("word", 123)
	assert.True(t, property.IsWord())
	assert.False(t, property.IsDate())
	assert.GreaterOrEqual(t, property.limit, uint(minLimit))
	assert.LessOrEqual(t, property.limit, uint(maxLimit))

	property = NewTaskProperty(time.Now().UTC().Format(dateFormatFromParam), 123)
	assert.True(t, property.IsDate())
	assert.False(t, property.IsWord())
	assert.GreaterOrEqual(t, property.limit, uint(minLimit))
	assert.LessOrEqual(t, property.limit, uint(maxLimit))
}
