package kv

import (
	"fmt"
	coreStorage "github.com/balerter/balerter/internal/core_storage"
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"testing"
)

type coreStorageMock struct {
	mock.Mock
	kv *coreStorageKVMock
}

func (m *coreStorageMock) KV() coreStorage.CoreStorageKV {
	return m.kv
}

func (m *coreStorageMock) Stop() error {
	return nil
}

func (m *coreStorageMock) Name() string {
	return ""
}

func (m *coreStorageMock) Alert() coreStorage.CoreStorageAlert {
	return nil
}

type coreStorageKVMock struct {
	mock.Mock
}

func (m *coreStorageKVMock) All() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *coreStorageKVMock) Get(_ string) (string, error) {
	return "", nil
}

func (m *coreStorageKVMock) Delete(_ string) error {
	return nil
}

func (m *coreStorageKVMock) Put(_, _ string) error {
	return nil
}

func (m *coreStorageKVMock) Upsert(_, _ string) error {
	return nil
}

func TestHandlerIndex(t *testing.T) {
	resultData := map[string]string{
		"f1": "v1",
		"f2": "v2",
	}

	am := &coreStorageMock{
		kv: &coreStorageKVMock{},
	}
	am.kv.On("All").Return(resultData, nil)

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

	am := &coreStorageMock{
		kv: &coreStorageKVMock{},
	}
	am.kv.On("All").Return(resultData, fmt.Errorf("error1"))

	f := HandlerIndex(am, zap.NewNop())

	rw := &httpTestify.TestResponseWriter{}
	req := &http.Request{URL: &url.URL{}}

	f(rw, req)

	assert.Equal(t, 500, rw.StatusCode)
	assert.Equal(t, "error1", rw.Header().Get("X-Error"))
	assert.Equal(t, "", rw.Output)
}
