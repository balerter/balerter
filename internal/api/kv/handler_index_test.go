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

	am := &coreStorage.KVMock{
		AllFunc: func() (map[string]string, error) {
			return resultData, nil
		},
	}

	kv := &KV{
		storage: am,
		logger:  zap.NewNop(),
	}

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	kv.handlerIndex(rw, req)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Contains(t, rw.Output, `"f1":"v1"`)
	assert.Contains(t, rw.Output, `"f2":"v2"`)
}

func TestHandlerIndex_ErrorGetFromStorage(t *testing.T) {
	am := &coreStorage.KVMock{
		AllFunc: func() (map[string]string, error) {
			return nil, fmt.Errorf("error1")
		},
	}

	kv := &KV{
		storage: am,
		logger:  zap.NewNop(),
	}

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	kv.handlerIndex(rw, req)

	assert.Equal(t, 500, rw.StatusCode)
	assert.Equal(t, "error1", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}
