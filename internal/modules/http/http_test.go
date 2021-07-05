package http

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestHttp_GetLoader(t *testing.T) {
	h := &HTTP{}

	loader := h.GetLoader(nil)

	L := lua.NewState()
	n := loader(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)
	assert.Equal(t, lua.LTFunction, v.RawGetString("request").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("get").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("post").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("put").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("delete").Type())

	assert.Equal(t, http.MethodGet, v.RawGetString("methodGet").String())
	assert.Equal(t, http.MethodHead, v.RawGetString("methodHead").String())
	assert.Equal(t, http.MethodPost, v.RawGetString("methodPost").String())
	assert.Equal(t, http.MethodPut, v.RawGetString("methodPut").String())
	assert.Equal(t, http.MethodPatch, v.RawGetString("methodPatch").String())
	assert.Equal(t, http.MethodDelete, v.RawGetString("methodDelete").String())
	assert.Equal(t, http.MethodConnect, v.RawGetString("methodConnect").String())
	assert.Equal(t, http.MethodOptions, v.RawGetString("methodOptions").String())
	assert.Equal(t, http.MethodTrace, v.RawGetString("methodTrace").String())
}

func TestModuleName(t *testing.T) {
	assert.Equal(t, "http", ModuleName())
}

func TestMethods(t *testing.T) {
	assert.Equal(t, []string{
		"request",
		"post",
		"get",
		"put",
		"delete",
	}, Methods())
}

func TestNew(t *testing.T) {
	h := New(zap.NewNop())
	assert.IsType(t, &HTTP{}, h)
}

func TestName(t *testing.T) {
	h := &HTTP{}
	assert.Equal(t, "http", h.Name())
}

func TestHTTP_Stop(t *testing.T) {
	hm := &httpClientMock{CloseIdleConnectionsFunc: func() {}}

	h := &HTTP{
		client: hm,
	}
	assert.NoError(t, h.Stop())

	s := hm.CloseIdleConnectionsCalls()

	assert.Equal(t, 1, len(s))
}
