package luaformatter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestValueToString(t *testing.T) {
	var s string
	var err error

	s, err = ValueToString(lua.LNumber(10))
	require.NoError(t, err)
	assert.Equal(t, "10", s)

	s, err = ValueToString(lua.LBool(true))
	require.NoError(t, err)
	assert.Equal(t, "true", s)

	s, err = ValueToString(lua.LString("foo"))
	require.NoError(t, err)
	assert.Equal(t, "foo", s)
}

func TestTableToString(t *testing.T) {
	t1 := &lua.LTable{}
	t1.RawSetString("f1", lua.LNumber(42))
	t1.RawSetString("f2", lua.LString("foo"))
	t1.RawSetString("f3", lua.LBool(true))
	t1.RawSetString("a", lua.LBool(false))

	t0 := &lua.LTable{}
	t0.RawSetString("f1", t1)
	t0.RawSetString("bar", lua.LString("baz"))
	t0.RawSetString("baz", lua.LNumber(100))

	s, err := TableToString(t0)
	require.NoError(t, err)
	assert.Equal(t, `{"bar":"baz","baz":100,"f1":{"a":false,"f1":42,"f2":"foo","f3":true}}`, s)
}

func TestTableToString_NilValue(t *testing.T) {
	t1 := &lua.LTable{}
	t1.RawSetString("foo", lua.LNil)
	s, err := TableToString(t1)
	require.NoError(t, err)
	assert.Equal(t, "{}", s)
}

func TestTableToString_BoolKey(t *testing.T) {
	t1 := &lua.LTable{}
	t1.RawSet(lua.LBool(true), lua.LNumber(42))
	_, err := TableToString(t1)
	require.Error(t, err)
	assert.Equal(t, "key must be a string", err.Error())
}

func TestTableToString_Errors(t *testing.T) {
	t1 := &lua.LTable{}
	t1.RawSet(lua.LNumber(42), lua.LNumber(42))
	_, err := TableToString(t1)
	require.Error(t, err)
	assert.Equal(t, "key must be a string", err.Error())

	_, err = TableToString(nil)
	require.Error(t, err)
	assert.Equal(t, "table is nil", err.Error())
}
