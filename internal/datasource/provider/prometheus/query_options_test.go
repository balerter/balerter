package prometheus

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_parseQueryOptions(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	tt := &lua.LTable{}
	tt.RawSetString("time", lua.LString("123"))
	luaState.Push(tt)
	opts, err := m.parseQueryOptions(luaState)
	require.NoError(t, err)
	assert.Equal(t, "123", opts.Time)
}

func Test_parseQueryOptions_without_options(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	tt := &lua.LTable{}
	//tt.RawSetString("time", lua.LString("123"))
	luaState.Push(tt)
	opts, err := m.parseQueryOptions(luaState)
	require.NoError(t, err)
	assert.Equal(t, "", opts.Time)
}

func Test_parseQueryOptions_options_not_a_table(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	luaState.Push(lua.LNumber(42))
	_, err := m.parseQueryOptions(luaState)
	require.Error(t, err)
	assert.Equal(t, "options must be a table", err.Error())
}

func Test_parseQueryOptions_option_time_not_a_string(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	tt := &lua.LTable{}
	tt.RawSetString("time", lua.LNumber(42))
	luaState.Push(tt)
	_, err := m.parseQueryOptions(luaState)
	require.Error(t, err)
	assert.Equal(t, "time must be a string", err.Error())
}
