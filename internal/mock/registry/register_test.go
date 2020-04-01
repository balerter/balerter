package registry

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_registry_register_no_calls(t *testing.T) {
	reg := New()

	ret := []lua.LValue{lua.LNumber(42), lua.LString("bar")}

	err := reg.Register("AnyValue", "foo", []lua.LValue{}, ret)
	require.NoError(t, err)

	e, ok := reg.responseEntries["foo"]
	require.True(t, ok)
	require.Equal(t, 1, len(e.responses))
	assert.Equal(t, ret, e.responses[0])
	assert.Equal(t, 0, len(e.entries))
}

func Test_registry_register_tree1(t *testing.T) {
	reg := New()

	ret := []lua.LValue{lua.LNumber(42), lua.LString("bar")}

	tbl := &lua.LTable{}
	tbl.RawSetString("foo", lua.LString("bar"))

	err := reg.Register("AnyValue", "foo", []lua.LValue{lua.LNumber(10), tbl}, ret)
	require.NoError(t, err)

	e, ok := reg.responseEntries["foo"]
	require.True(t, ok)
	require.Equal(t, 0, len(e.responses))
	assert.Equal(t, 1, len(e.entries))

	e1, ok := e.entries["number@10"]
	require.True(t, ok)
	assert.Equal(t, 0, len(e1.responses))
	assert.Equal(t, 1, len(e1.entries))

	e2, ok := e1.entries["table@{\"foo\":\"bar\"}"]
	require.True(t, ok)
	assert.Equal(t, 1, len(e2.responses))
	assert.Equal(t, 0, len(e2.entries))
	assert.Equal(t, ret, e2.responses[0])
}

func Test_registry_register_any_value(t *testing.T) {
	reg := New()

	ret := []lua.LValue{lua.LNumber(42), lua.LString("bar")}

	tbl := &lua.LTable{}
	tbl.RawSetString("foo", lua.LString("bar"))

	err := reg.Register("AnyValue", "foo", []lua.LValue{lua.LNumber(10), lua.LString("AnyValue")}, ret)
	require.NoError(t, err)

	e, ok := reg.responseEntries["foo"]
	require.True(t, ok)
	require.Equal(t, 0, len(e.responses))
	assert.Equal(t, 1, len(e.entries))

	e1, ok := e.entries["number@10"]
	require.True(t, ok)
	assert.Equal(t, 0, len(e1.responses))
	assert.Equal(t, 1, len(e1.entries))

	e2, ok := e1.entries["AnyValue"]
	require.True(t, ok)
	assert.Equal(t, 1, len(e2.responses))
	assert.Equal(t, 0, len(e2.entries))
	assert.Equal(t, ret, e2.responses[0])
}

func Test_registry_register_multiple_rets(t *testing.T) {
	reg := New()

	ret1 := []lua.LValue{lua.LNumber(42), lua.LString("bar")}
	ret2 := []lua.LValue{lua.LNumber(10), lua.LString("baz")}

	err := reg.Register("AnyValue", "foo", []lua.LValue{}, ret1)
	require.NoError(t, err)
	err = reg.Register("AnyValue", "foo", []lua.LValue{}, ret2)
	require.NoError(t, err)

	e, ok := reg.responseEntries["foo"]
	require.True(t, ok)
	assert.Equal(t, 0, len(e.entries))
	require.Equal(t, 2, len(e.responses))
	assert.Equal(t, ret1, e.responses[0])
	assert.Equal(t, ret2, e.responses[1])
}
