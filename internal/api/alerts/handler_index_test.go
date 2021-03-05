package alerts

import (
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlerIndex_bad_levels(t *testing.T) {
	a := &Alerts{
		logger: zap.NewNop(),
	}

	rw := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	req.URL.RawQuery = "levels=success,foo"

	a.handlerIndex(rw, req)

	assert.Equal(t, 400, rw.Code)
	assert.Equal(t, "error parse level foo, bad level\n", rw.Body.String())
}

func TestHandlerIndex_bad_get_index(t *testing.T) {
	m := &coreStorageAlertMock{}

	a := &Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

	m.On("Index", mock.Anything).Return(nil, fmt.Errorf("err1"))

	rw := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	a.handlerIndex(rw, req)

	assert.Equal(t, 500, rw.Code)
	assert.Equal(t, "internal error\n", rw.Body.String())
}

func TestHandlerIndex(t *testing.T) {
	m := &coreStorageAlertMock{}

	a := &Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

	al := &alert2.Alert{
		Name:       "1",
		Level:      2,
		LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
		Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
		Count:      3,
	}

	m.On("Index", mock.Anything).Return(alert2.Alerts{al}, nil)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	a.handlerIndex(rw, req)

	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, `[{"name":"1","level":"warning","level_num":2,"count":3,"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"}]`, rw.Body.String())
}
