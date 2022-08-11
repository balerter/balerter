package alerts

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandlerGet_empty_name(t *testing.T) {
	a := Alerts{}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	a.handlerGet(rw, req)

	assert.Equal(t, "empty name\n", rw.Body.String())
	assert.Equal(t, 400, rw.Code)
}

func TestHandlerGet_get_error(t *testing.T) {
	m := &corestorage.AlertMock{
		GetFunc: func(name string) (*alert.Alert, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

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
	m := &corestorage.AlertMock{
		GetFunc: func(name string) (*alert.Alert, error) {
			return nil, nil
		},
	}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

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
	al := &alert.Alert{
		Name:       "1",
		Level:      2,
		LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
		Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
		Count:      3,
	}

	m := &corestorage.AlertMock{
		GetFunc: func(name string) (*alert.Alert, error) {
			return al, nil
		},
	}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

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
