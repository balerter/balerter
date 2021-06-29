package alerts

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type coreStorageAlertMock struct {
	mock.Mock
}

func (m *coreStorageAlertMock) Update(name string, level alert.Level) (*alert.Alert, bool, error) {
	args := m.Called(name, level)
	a := args.Get(0)
	if a == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return a.(*alert.Alert), args.Bool(1), args.Error(2)
}

func (m *coreStorageAlertMock) Index(levels []alert.Level) (alert.Alerts, error) {
	args := m.Called(levels)
	a := args.Get(0)
	if a == nil {
		return nil, args.Error(1)
	}
	return a.(alert.Alerts), args.Error(1)
}

func (m *coreStorageAlertMock) Get(name string) (*alert.Alert, error) {
	args := m.Called(name)
	a := args.Get(0)
	if a == nil {
		return nil, args.Error(1)
	}
	return a.(*alert.Alert), args.Error(1)
}

func TestHandlerGet_empty_name(t *testing.T) {
	a := Alerts{}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	a.handlerGet(rw, req)

	assert.Equal(t, "empty name\n", rw.Body.String())
	assert.Equal(t, 400, rw.Code)
}

func TestHandlerGet_get_error(t *testing.T) {
	m := &coreStorageAlertMock{}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

	m.On("Get", mock.Anything).Return(nil, fmt.Errorf("err1"))

	rw := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(t, err)

	a.handlerGet(rw, req)

	assert.Equal(t, "error get alert\n", rw.Body.String())
	assert.Equal(t, 500, rw.Code)
}

func TestHandlerGet_alert_not_found(t *testing.T) {
	m := &coreStorageAlertMock{}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

	m.On("Get", mock.Anything).Return(nil, nil)

	rw := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(t, err)

	a.handlerGet(rw, req)

	assert.Equal(t, "alert not found\n", rw.Body.String())
	assert.Equal(t, 404, rw.Code)
}

func TestHandlerGet(t *testing.T) {
	m := &coreStorageAlertMock{}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

	al := &alert.Alert{
		Name:       "1",
		Level:      2,
		LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
		Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
		Count:      3,
	}

	m.On("Get", mock.Anything).Return(al, nil)

	rw := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(t, err)

	a.handlerGet(rw, req)

	assert.Equal(t, `{"name":"1","level":"warning","level_num":2,"count":3,`+
		`"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"}`, rw.Body.String())
	assert.Equal(t, 200, rw.Code)
}
