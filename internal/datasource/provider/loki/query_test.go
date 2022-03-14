package loki

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func Test_doQuery_getQuery_no_query(t *testing.T) {
	m := &Loki{}

	luaState := lua.NewState()

	_, err := m.getQuery(luaState)
	require.Error(t, err)
	assert.Equal(t, "query must be not empty", err.Error())
}

func Test_doQuery_getQuery_empty_query(t *testing.T) {
	m := &Loki{}

	luaState := lua.NewState()
	luaState.Push(lua.LString(""))

	_, err := m.getQuery(luaState)
	require.Error(t, err)
	assert.Equal(t, "query must be not empty", err.Error())
}

func Test_doQuery_error_parse_options(t *testing.T) {
	m := &Loki{
		url:    &url.URL{Host: "domain.com"},
		logger: zap.NewNop(),
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString("q1"))
	luaState.Push(lua.LString("opts"))

	n := m.doQuery(luaState)

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTNil, luaState.Get(3).Type())
	assert.Equal(t, lua.LTString, luaState.Get(4).Type())
	assert.Equal(t, "error parse query options", luaState.Get(4).String())
}

func Test_doQuery_error_send_request(t *testing.T) {
	mm := &httpClientMock{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	m := &Loki{
		logger: zap.NewNop(),
		url:    &url.URL{Host: "domain.com"},
		client: mm,
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString("q1"))

	n := m.doQuery(luaState)

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTNil, luaState.Get(2).Type())
	assert.Equal(t, lua.LTString, luaState.Get(3).Type())
	assert.Equal(t, "error send query to loki: err1", luaState.Get(3).String())
}

func Test_doRange_error_parse_options(t *testing.T) {
	m := &Loki{
		url:    &url.URL{Host: "domain.com"},
		logger: zap.NewNop(),
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString("q1"))
	luaState.Push(lua.LString("opts"))

	n := m.doRange(luaState)

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTNil, luaState.Get(3).Type())
	assert.Equal(t, lua.LTString, luaState.Get(4).Type())
	assert.Equal(t, "error parse range options", luaState.Get(4).String())
}

func Test_doRange_error_send_request(t *testing.T) {
	mm := &httpClientMock{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	m := &Loki{
		logger: zap.NewNop(),
		url:    &url.URL{Host: "domain.com"},
		client: mm,
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString("q1"))

	n := m.doRange(luaState)

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTNil, luaState.Get(2).Type())
	assert.Equal(t, lua.LTString, luaState.Get(3).Type())
	assert.Equal(t, "error send query to loki: err1", luaState.Get(3).String())
}

func Test_do_unexpected_model_type(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"data":{"resultType": "vector","result":[]}}`))),
	}

	mm := &httpClientMock{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return resp, nil
		},
	}

	m := &Loki{
		logger: zap.NewNop(),
		client: mm,
	}

	luaState := lua.NewState()

	n := m.do(luaState, "")

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTNil, luaState.Get(1).Type())
	assert.Equal(t, lua.LTString, luaState.Get(2).Type())
	assert.Equal(t, "query error: unexpected loki model type: vector", luaState.Get(2).String())
}

func Test_doQuery(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"status": "success","data":{"resultType": "streams","result": [{"stream": {},"values": []}]}}`))),
	}
	mm := &httpClientMock{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return resp, nil
		},
	}
	m := &Loki{
		logger: zap.NewNop(),
		client: mm,
	}

	luaState := lua.NewState()

	n := m.do(luaState, "")

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTTable, luaState.Get(1).Type())
	assert.Equal(t, lua.LTNil, luaState.Get(2).Type())
}

func Test_doQuery_badQuery(t *testing.T) {
	m := &Loki{
		logger: zap.NewNop(),
	}

	luaState := lua.NewState()

	n := m.doQuery(luaState)

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTNil, luaState.Get(1).Type())
	assert.Equal(t, lua.LTString, luaState.Get(2).Type())
	assert.Equal(t, "query must be not empty", luaState.Get(2).String())
}

func Test_doRange_badQuery(t *testing.T) {
	m := &Loki{
		logger: zap.NewNop(),
	}

	luaState := lua.NewState()

	n := m.doRange(luaState)

	assert.Equal(t, 2, n)
	assert.Equal(t, lua.LTNil, luaState.Get(1).Type())
	assert.Equal(t, lua.LTString, luaState.Get(2).Type())
	assert.Equal(t, "query must be not empty", luaState.Get(2).String())
}
