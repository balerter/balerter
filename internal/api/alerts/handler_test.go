package alerts

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type coreStorageMock struct {
	mock.Mock
}

func (m *coreStorageMock) GetOrNew(string) (*alert.Alert, error) {
	args := m.Called()
	return args.Get(0).(*alert.Alert), args.Error(1)
}
func (m *coreStorageMock) All() ([]*alert.Alert, error) {
	args := m.Called()
	return args.Get(0).([]*alert.Alert), args.Error(1)
}
func (m *coreStorageMock) Release(_ *alert.Alert) {
}

func TestHandler_ErrorGetAlerts(t *testing.T) {
	var resultData []*alert.Alert

	am := &coreStorageMock{}
	am.On("All").Return(resultData, fmt.Errorf("error1"))

	f := Handler(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	f(rw, req)

	assert.Equal(t, 500, rw.StatusCode)
	assert.Equal(t, "error1", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}

func TestHandler(t *testing.T) {
	var resultData []*alert.Alert

	a1 := alert.AcquireAlert()
	a1.SetName("foo")
	a1.UpdateLevel(alert.LevelError)
	a1.Inc()
	resultData = append(resultData, a1)

	updatedAt := a1.GetLastChangeTime().Format(time.RFC3339)

	am := &coreStorageMock{}
	am.On("All").Return(resultData, nil)

	f := Handler(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	f(rw, req)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Equal(t, `[{"name":"foo","level":"error","count":1,"updated_at":"`+updatedAt+`"}]`, rw.Output)
}

func TestHandler_BadLevelArgument(t *testing.T) {
	var resultData []*alert.Alert

	am := &coreStorageMock{}
	am.On("All").Return(resultData, nil)

	f := Handler(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{RawQuery: "level=foo"}}

	f(rw, req)

	assert.Equal(t, 400, rw.StatusCode)
	assert.Equal(t, "bad level value", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}
