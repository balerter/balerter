package runtime

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type runnerMock struct {
	mock.Mock
}

func (r *runnerMock) RunScript(name string, req *http.Request) error {
	args := r.Called(name, req)
	return args.Error(0)
}

func Test_handlerRun_empty_name(t *testing.T) {
	rt := &Runtime{}

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rt.handlerRun(rw, req)

	assert.Equal(t, "empty name\n", rw.Body.String())
	assert.Equal(t, 400, rw.Code)
}

func Test_handlerRun_runScriptError(t *testing.T) {
	m := &runnerMock{}

	rt := &Runtime{
		runner: m,
	}

	m.On("RunScript", mock.Anything, mock.Anything).Return(fmt.Errorf("err1"))

	rw := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(t, err)

	rt.handlerRun(rw, req)

	assert.Equal(t, "err1\n", rw.Body.String())
	assert.Equal(t, 400, rw.Code)
}

func Test_handlerRun(t *testing.T) {
	m := &runnerMock{}

	rt := &Runtime{
		runner: m,
	}

	m.On("RunScript", mock.Anything, mock.Anything).Return(nil)

	rw := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("name", "foo")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(t, err)

	rt.handlerRun(rw, req)

	assert.Equal(t, 200, rw.Code)
}
