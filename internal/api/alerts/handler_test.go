package alerts

import (
	"github.com/balerter/balerter/internal/alert/alert"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type alertManagerAPIerMock struct {
	mock.Mock
}

func (m *alertManagerAPIerMock) GetAlerts() []*alertManager.AlertInfo {
	return m.Called().Get(0).([]*alertManager.AlertInfo)
}

func TestHandler(t *testing.T) {
	resultData := []*alertManager.AlertInfo{
		{Name: "foo", Level: alert.LevelError, Count: 5, LastChange: time.Date(2020, 01, 01, 10, 10, 10, 0, time.UTC)},
	}

	am := &alertManagerAPIerMock{}
	am.On("GetAlerts").Return(resultData)

	f := Handler(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	f(rw, req)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Equal(t, `[{"name":"foo","level":"error","count":5,"updated_at":"2020-01-01T10:10:10Z"}]`, rw.Output)
}

func TestHandler_BadLevelArgument(t *testing.T) {
	resultData := []*alertManager.AlertInfo{
		{Name: "foo", Level: alert.LevelError},
	}

	am := &alertManagerAPIerMock{}
	am.On("GetAlerts").Return(resultData)

	f := Handler(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{RawQuery: "level=foo"}}

	f(rw, req)

	assert.Equal(t, 400, rw.StatusCode)
	assert.Equal(t, "bad level value", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}
