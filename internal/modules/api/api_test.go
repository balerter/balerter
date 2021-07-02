package api

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"net/http"
	"net/url"
	"testing"
)

func TestModuleName(t *testing.T) {
	assert.Equal(t, "api", ModuleName())
}

func TestMethods(t *testing.T) {
	assert.Equal(t, []string{"is_api",
		"query",
		"url",
		"body",
		"host",
		"method",
	}, Methods())
}

func TestNew(t *testing.T) {
	p := New()
	assert.IsType(t, &API{}, p)
}

type badReader struct {
	mock.Mock
}

func (m *badReader) Read(b []byte) (int, error) {
	a := m.Called(b)
	return a.Int(0), a.Error(1)
}

func TestAPI_FillData_error_read_body(t *testing.T) {
	a := &API{}

	r := &badReader{}
	r.On("Read", mock.Anything).Return(0, fmt.Errorf("err1"))

	req, err := http.NewRequest("POST", "", r)
	require.NoError(t, err)
	req.URL = &url.URL{}

	err = a.FillData(req)
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestAPI_FillData(t *testing.T) {
	a := &API{query: map[string][]string{}}

	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte("foo")))
	require.NoError(t, err)
	req.URL = &url.URL{
		RawQuery: "foo=1&bar=2",
	}

	err = a.FillData(req)
	require.NoError(t, err)

	assert.Equal(t, 2, len(a.query))

	v, ok := a.query["foo"]
	assert.True(t, ok)
	assert.Equal(t, []string{"1"}, v)

	v, ok = a.query["bar"]
	assert.True(t, ok)
	assert.Equal(t, []string{"2"}, v)
}

func TestAPI_GetLoader(t *testing.T) {
	a := &API{}

	f := a.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func TestAPI_Stop(t *testing.T) {
	a := &API{}
	assert.NoError(t, a.Stop())
}

func Test_funcIsAPI(t *testing.T) {
	a := &API{
		isAPI: true,
	}
	luaState := lua.NewState()
	n := a.funcIsAPI(luaState)
	assert.Equal(t, 1, n)
	assert.Equal(t, "true", luaState.Get(1).String())
}

func Test_funcQuery(t *testing.T) {
	a := &API{
		query: map[string][]string{"foo": {"bar"}},
	}
	luaState := lua.NewState()
	n := a.funcQuery(luaState)
	assert.Equal(t, 1, n)

	assert.Equal(t, lua.LTTable, luaState.Get(1).Type())
}

func Test_funcString(t *testing.T) {
	a := &API{}
	f := a.funcString("foo")
	luaState := lua.NewState()
	n := f(luaState)
	assert.Equal(t, 1, n)
	assert.Equal(t, "foo", luaState.Get(1).String())
}
