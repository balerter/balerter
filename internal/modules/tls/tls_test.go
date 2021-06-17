package tls

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestModuleName(t *testing.T) {
	assert.Equal(t, "tls", ModuleName())
}

func TestMethods(t *testing.T) {
	assert.Equal(t, []string{"get"}, Methods())
}

func TestNew(t *testing.T) {
	m := New()
	assert.IsType(t, &TLS{}, m)
}

func TestStop(t *testing.T) {
	m := &TLS{}
	assert.Nil(t, m.Stop())
}

func TestGetLoader(t *testing.T) {
	m := &TLS{}

	loader := m.GetLoader(nil)

	luaState := lua.NewState()

	n := loader(luaState)
	assert.Equal(t, 1, n)

	v := luaState.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func Test_get_wrong_param(t *testing.T) {
	m := &TLS{}

	luaState := lua.NewState()

	n := m.get(luaState)
	assert.Equal(t, 2, n)

	assert.Equal(t, lua.LTNil, luaState.Get(1).Type())

	assert.Equal(t, lua.LTString, luaState.Get(2).Type())
	assert.Equal(t, "parameter must be a string", luaState.Get(2).String())
}

func Test_get_wrong_param_2(t *testing.T) {
	m := &TLS{}

	luaState := lua.NewState()
	luaState.Push(lua.LNumber(42))

	n := m.get(luaState)
	assert.Equal(t, 2, n)

	assert.Equal(t, lua.LTNil, luaState.Get(2).Type())

	assert.Equal(t, lua.LTString, luaState.Get(3).Type())
	assert.Equal(t, "parameter must be a string", luaState.Get(3).String())
}

func Test_get_dial_error(t *testing.T) {
	dialFunc := func(network, addr string, config *tls.Config) (*tls.Conn, error) {
		return nil, fmt.Errorf("err1")
	}

	m := &TLS{
		dialFunc: dialFunc,
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString("domain.com"))

	n := m.get(luaState)
	assert.Equal(t, 2, n)

	assert.Equal(t, lua.LTNil, luaState.Get(2).Type())

	assert.Equal(t, lua.LTString, luaState.Get(3).Type())
	assert.Equal(t, "error dial to host domain.com, err1", luaState.Get(3).String())
}
