package alerts

import (
	"bytes"
	"context"
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"
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

type chManagerMock struct {
	mock.Mock
}

func (m *chManagerMock) Send(a *alert2.Alert, text string, options *alert2.Options) {
	m.Called(a, text, options)
}

func TestHandlerUpdate_empty_name(t *testing.T) {
	a := Alerts{}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	a.handlerUpdate(rw, req)

	assert.Equal(t, "empty name\n", rw.Body.String())
	assert.Equal(t, 400, rw.Code)
}

type readCloserMock struct {
	mock.Mock
}

func (r *readCloserMock) Read(p []byte) (int, error) {
	args := r.Called(p)
	return args.Int(0), args.Error(1)
}

func (r *readCloserMock) Close() error {
	args := r.Called()
	return args.Error(0)
}

func TestHandlerUpdate_error_read_body(t *testing.T) {
	a := Alerts{
		logger: zap.NewNop(),
	}

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	mr := &readCloserMock{}
	mr.On("Close").Return(nil)
	mr.On("Read", mock.Anything).Return(0, fmt.Errorf("err1"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", mr)
	require.NoError(t, err)

	rw := httptest.NewRecorder()

	a.handlerUpdate(rw, req)

	assert.Equal(t, "error read body\n", rw.Body.String())
	assert.Equal(t, 500, rw.Code)
}

func TestHandlerUpdate_error_level_from_string(t *testing.T) {
	a := Alerts{
		logger: zap.NewNop(),
	}

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	mr := bytes.NewBuffer([]byte(`{}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", mr)
	require.NoError(t, err)

	rw := httptest.NewRecorder()

	a.handlerUpdate(rw, req)

	assert.Equal(t, "error parse level , bad level\n", rw.Body.String())
	assert.Equal(t, 400, rw.Code)
}

func TestHandlerUpdate_error_update(t *testing.T) {
	m := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return nil, false, fmt.Errorf("err1")
		},
	}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	mr := bytes.NewBuffer([]byte(`{"level":"success","text":""}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", mr)
	require.NoError(t, err)

	rw := httptest.NewRecorder()

	a.handlerUpdate(rw, req)

	assert.Equal(t, "error update alert\n", rw.Body.String())
	assert.Equal(t, 500, rw.Code)
}

func TestHandlerUpdate_level_was_updated(t *testing.T) {
	al := &alert2.Alert{
		Name:       "1",
		Level:      2,
		LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
		Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
		Count:      3,
	}

	m := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return al, true, nil
		},
	}

	ch := &chManagerMock{}

	a := Alerts{
		alertManager: m,
		chManager:    ch,
		logger:       zap.NewNop(),
	}

	ch.On("Send", mock.Anything, mock.Anything, mock.Anything).Return()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	mr := bytes.NewBuffer([]byte(`{"level":"success","text":""}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", mr)
	require.NoError(t, err)

	rw := httptest.NewRecorder()

	a.handlerUpdate(rw, req)

	ch.AssertCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything)

	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, `{"name":"1","level":"warning","level_num":2,"count":3,`+
		`"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"}`, rw.Body.String())

	ch.AssertExpectations(t)
}

func TestHandlerUpdate_level_was_not_updated(t *testing.T) {
	al := &alert2.Alert{
		Name:       "1",
		Level:      2,
		LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
		Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
		Count:      3,
	}

	m := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return al, false, nil
		},
	}

	a := Alerts{
		alertManager: m,
		logger:       zap.NewNop(),
	}

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	mr := bytes.NewBuffer([]byte(`{"level":"success","text":""}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", mr)
	require.NoError(t, err)

	rw := httptest.NewRecorder()

	a.handlerUpdate(rw, req)

	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, `{"name":"1","level":"warning","level_num":2,"count":3,`+
		`"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"}`, rw.Body.String())
}

func TestHandlerUpdate_resend(t *testing.T) {
	al := &alert2.Alert{
		Name:       "1",
		Level:      2,
		LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
		Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
		Count:      10,
	}

	m := &corestorage.AlertMock{
		UpdateFunc: func(name string, level alert2.Level) (*alert2.Alert, bool, error) {
			return al, false, nil
		},
	}

	ch := &chManagerMock{}

	a := Alerts{
		alertManager: m,
		chManager:    ch,
		logger:       zap.NewNop(),
	}

	ch.On("Send", mock.Anything, mock.Anything, mock.Anything).Return()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)

	mr := bytes.NewBuffer([]byte(`{"level":"success","text":"","repeat":2}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", mr)
	require.NoError(t, err)

	rw := httptest.NewRecorder()

	a.handlerUpdate(rw, req)

	ch.AssertCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything)

	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, `{"name":"1","level":"warning","level_num":2,"count":10,`+
		`"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"}`, rw.Body.String())

	ch.AssertExpectations(t)
}
