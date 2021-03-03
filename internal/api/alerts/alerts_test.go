package alerts

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAlerts_Handler(t *testing.T) {
	am := &Alerts{}

	r := &chiMock{}
	r.On("Get", "/", mock.AnythingOfType("http.HandlerFunc"))
	r.On("Get", "/{name}", mock.AnythingOfType("http.HandlerFunc"))
	r.On("Post", "/{name}", mock.AnythingOfType("http.HandlerFunc"))

	am.Handler(r)

	r.AssertCalled(t, "Get", "/", mock.AnythingOfType("http.HandlerFunc"))
	r.AssertCalled(t, "Get", "/{name}", mock.AnythingOfType("http.HandlerFunc"))
	r.AssertCalled(t, "Post", "/{name}", mock.AnythingOfType("http.HandlerFunc"))

	r.AssertExpectations(t)
}

func TestNew(t *testing.T) {
	am := New(nil, nil, nil)
	assert.IsType(t, &Alerts{}, am)
}
