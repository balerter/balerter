package prometheus

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_parseRangeOptions(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	tt := &lua.LTable{}
	tt.RawSetString("start", lua.LString("1"))
	tt.RawSetString("end", lua.LString("2"))
	tt.RawSetString("step", lua.LString("3"))
	luaState.Push(tt)
	opts, err := m.parseRangeOptions(luaState)
	require.NoError(t, err)
	assert.Equal(t, "1", opts.Start)
	assert.Equal(t, "2", opts.End)
	assert.Equal(t, "3", opts.Step)
}

func Test_parseRangeOptions_options_not_a_table(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	luaState.Push(lua.LNumber(42))
	_, err := m.parseRangeOptions(luaState)
	require.Error(t, err)
	assert.Equal(t, "options must be a table", err.Error())
}

func Test_parseRangeOptions_start_not_a_string(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	tt := &lua.LTable{}
	tt.RawSetString("start", lua.LNumber(1))
	tt.RawSetString("end", lua.LString("2"))
	tt.RawSetString("step", lua.LString("3"))
	luaState.Push(tt)
	_, err := m.parseRangeOptions(luaState)
	require.Error(t, err)
	assert.Equal(t, "start must be a string", err.Error())
}

func Test_parseRangeOptions_end_not_a_string(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	tt := &lua.LTable{}
	tt.RawSetString("start", lua.LString("1"))
	tt.RawSetString("end", lua.LNumber(2))
	tt.RawSetString("step", lua.LString("3"))
	luaState.Push(tt)
	_, err := m.parseRangeOptions(luaState)
	require.Error(t, err)
	assert.Equal(t, "end must be a string", err.Error())
}

func Test_parseRangeOptions_step_not_a_string(t *testing.T) {
	m := &Prometheus{}
	luaState := lua.NewState()
	luaState.Push(lua.LNil)
	tt := &lua.LTable{}
	tt.RawSetString("start", lua.LString("1"))
	tt.RawSetString("end", lua.LString("2"))
	tt.RawSetString("step", lua.LNumber(3))
	luaState.Push(tt)
	_, err := m.parseRangeOptions(luaState)
	require.Error(t, err)
	assert.Equal(t, "step must be a string", err.Error())
}
