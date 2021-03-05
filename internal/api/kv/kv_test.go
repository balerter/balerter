package kv

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNew(t *testing.T) {
	kv := New(nil, nil)
	assert.IsType(t, &KV{}, kv)
}

func TestHandler(t *testing.T) {
	am := &KV{}

	r := &chiMock{}
	r.On("Get", "/", mock.AnythingOfType("http.HandlerFunc"))

	am.Handler(r)

	r.AssertCalled(t, "Get", "/", mock.AnythingOfType("http.HandlerFunc"))

	r.AssertExpectations(t)
}
