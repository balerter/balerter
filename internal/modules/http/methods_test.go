package http

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"io"
	"net/http"
	"testing"
)

func TestHTTP_send_error_parseRequestArgs(t *testing.T) {
	h := &HTTP{}

	f := h.send("POST")

	ls := lua.NewState()

	n := f(ls)

	assert.Equal(t, 2, n)
	assert.Equal(t, "error parse args, uri argument must be a string", ls.Get(2).String())
}

func TestHTTP_send_error_send_request(t *testing.T) {
	hc := &httpClientMock{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	h := &HTTP{
		logger: zap.NewNop(),
		client: hc,
	}

	f := h.send("POST")

	ls := lua.NewState()
	ls.Push(lua.LString("localhost"))
	ls.Push(lua.LString("body"))

	n := f(ls)

	assert.Equal(t, 2, n)
	assert.Equal(t, "error send request, err1", ls.Get(4).String())
}

func TestHTTP_send(t *testing.T) {
	hc := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			r := &http.Response{}
			r.Body = io.NopCloser(bytes.NewBuffer([]byte("response")))
			return r, nil
		},
	}

	h := &HTTP{
		logger: zap.NewNop(),
		client: hc,
	}

	f := h.send("POST")

	ls := lua.NewState()
	ls.Push(lua.LString("localhost"))
	ls.Push(lua.LString("body"))

	n := f(ls)

	assert.Equal(t, 1, n)
	tbl := ls.Get(3)
	require.Equal(t, lua.LTTable, tbl.Type())
}

func TestHTTP_request_error_argument(t *testing.T) {
	h := &HTTP{}

	ls := lua.NewState()

	n := h.request(ls)

	assert.Equal(t, 2, n)
	require.Equal(t, "argument must be a table", ls.Get(2).String())
}

func TestHTTP_request_error_parse_args(t *testing.T) {
	h := &HTTP{
		logger: zap.NewNop(),
	}

	ls := lua.NewState()
	args := &lua.LTable{}
	args.RawSetString("method", lua.LString("foo"))
	ls.Push(args)

	n := h.request(ls)

	assert.Equal(t, 2, n)
	require.Equal(t, "error parse arguments, bad http method foo", ls.Get(3).String())
}

func TestHTTP_request_error_send_request(t *testing.T) {
	hc := &httpClientMock{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	h := &HTTP{
		logger: zap.NewNop(),
		client: hc,
	}

	ls := lua.NewState()
	ls.Push(&lua.LTable{})

	n := h.request(ls)

	assert.Equal(t, 2, n)
	require.Equal(t, "error send request, err1", ls.Get(3).String())
}

func TestHTTP_request(t *testing.T) {
	hc := &httpClientMock{
		DoFunc: func(_ *http.Request) (*http.Response, error) {
			r := &http.Response{}
			r.Body = io.NopCloser(bytes.NewBuffer([]byte("foo")))
			return r, nil
		},
	}

	h := &HTTP{
		logger: zap.NewNop(),
		client: hc,
	}

	ls := lua.NewState()
	ls.Push(&lua.LTable{})

	n := h.request(ls)

	assert.Equal(t, 1, n)
	assert.Equal(t, lua.LTTable, ls.Get(2).Type())
}
