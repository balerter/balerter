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
