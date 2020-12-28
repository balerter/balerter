package alerts

import (
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHandler_ErrorGetAlerts(t *testing.T) {
	var resultData []*alert2.Alert

	am := coreStorage.NewMock("")
	am.AlertMock().On("All").Return(resultData, fmt.Errorf("error1"))

	a := &Alerts{
		storage: am.Alert(),
		logger:  zap.NewNop(),
	}

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	a.handlerIndex(rw, req)

	assert.Equal(t, 500, rw.StatusCode)
	assert.Equal(t, "error1", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}

func TestHandler(t *testing.T) {
	var resultData []*alert2.Alert

	a1 := alert2.AcquireAlert()
	a1.SetName("foo")
	a1.UpdateLevel(alert2.LevelError)
	a1.Inc()
	resultData = append(resultData, a1)

	updatedAt := a1.GetLastChangeTime().Format(time.RFC3339)

	am := coreStorage.NewMock("")
	am.AlertMock().On("All").Return(resultData, nil)

	a := &Alerts{
		storage: am.Alert(),
		logger:  zap.NewNop(),
	}

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	a.handlerIndex(rw, req)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Equal(t, `[{"name":"foo","level":"error","count":1,"updated_at":"`+updatedAt+`"}]`, rw.Output)
}

func TestHandler_BadLevelArgument(t *testing.T) {
	var resultData []*alert2.Alert

	am := coreStorage.NewMock("")
	am.AlertMock().On("All").Return(resultData, nil)

	a := &Alerts{
		storage: am.Alert(),
		logger:  zap.NewNop(),
	}

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{RawQuery: "level=foo"}}

	a.handlerIndex(rw, req)

	assert.Equal(t, 400, rw.StatusCode)
	assert.Equal(t, "bad level value", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}
