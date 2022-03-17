package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPhotoMessage(t *testing.T) {
	m := NewPhotoMessage(1, "a", "b")
	assert.IsType(t, &PhotoMessage{}, m)
	assert.Equal(t, int64(1), m.ChatID)
	assert.Equal(t, "a", m.Photo)
	assert.Equal(t, "b", m.Caption)
}

func TestNewTextMessage(t *testing.T) {
	m := NewTextMessage(1, "a")
	assert.IsType(t, &TextMessage{}, m)
	assert.Equal(t, int64(1), m.ChatID)
	assert.Equal(t, "a", m.Text)
}
