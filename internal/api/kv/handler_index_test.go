package kv

import (
	"fmt"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"testing"
)

func TestHandlerIndex(t *testing.T) {
	resultData := map[string]string{
		"f1": "v1",
		"f2": "v2",
	}

	am := coreStorage.NewMock("")
	am.KVMock().On("All").Return(resultData, nil)

	f := HandlerIndex(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	f(rw, req)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Contains(t, rw.Output, `{"name":"f1","value":"v1"}`)
	assert.Contains(t, rw.Output, `{"name":"f2","value":"v2"}`)
}

func TestHandlerIndex_ErrorGetFromStorage(t *testing.T) {
	resultData := map[string]string{}

	am := coreStorage.NewMock("")
	am.KVMock().On("All").Return(resultData, fmt.Errorf("error1"))

	f := HandlerIndex(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	f(rw, req)

	assert.Equal(t, 500, rw.StatusCode)
	assert.Equal(t, "error1", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}
