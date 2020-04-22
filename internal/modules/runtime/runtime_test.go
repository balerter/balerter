package runtime

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestNew(t *testing.T) {
	m := New("1", true, true, "2", "3", nil)

	assert.IsType(t, &Runtime{}, m)
	assert.Equal(t, "1", m.logLevel)
	assert.Equal(t, true, m.isDebug)
	assert.Equal(t, true, m.isOnce)
	assert.Equal(t, "2", m.withScript)
	assert.Equal(t, "3", m.configSource)
}

func TestName(t *testing.T) {
	m := &Runtime{}

	assert.Equal(t, "runtime", m.Name())
}

func TestGetLoader(t *testing.T) {
	m := &Runtime{}

	f := m.GetLoader(nil)

	L := lua.NewState()

	n := f(L)

	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.Equal(t, lua.LTFunction, v.RawGetString(method).Type())
	}
}

func TestStop(t *testing.T) {
	m := &Runtime{}

	assert.NoError(t, m.Stop())
}

func Test_returnBool(t *testing.T) {
	m := &Runtime{}

	L := lua.NewState()
	f := m.returnBool(true)
	n := f(L)
	assert.Equal(t, 1, n)
	v := L.Get(1)
	assert.Equal(t, lua.LTBool, v.Type())
	assert.Equal(t, "true", v.String())

	L = lua.NewState()
	f = m.returnBool(false)
	n = f(L)
	assert.Equal(t, 1, n)
	v = L.Get(1)
	assert.Equal(t, lua.LTBool, v.Type())
	assert.Equal(t, "false", v.String())
}

func Test_returnString(t *testing.T) {
	m := &Runtime{}

	L := lua.NewState()
	f := m.returnString("foo")
	n := f(L)
	assert.Equal(t, 1, n)
	v := L.Get(1)
	assert.Equal(t, lua.LTString, v.Type())
	assert.Equal(t, "foo", v.String())
}
