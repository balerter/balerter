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
