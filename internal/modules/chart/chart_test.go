package chart

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestModuleName(t *testing.T) {
	assert.Equal(t, "chart", ModuleName())
}

func TestMethods(t *testing.T) {
	m := Methods()
	require.Equal(t, 1, len(m))
	assert.Equal(t, "render", m[0])
}

func TestNew(t *testing.T) {
	ch := New(nil)
	assert.IsType(t, &Chart{}, ch)
}

func TestChart_Stop(t *testing.T) {
	ch := &Chart{}
	assert.NoError(t, ch.Stop())
}

func TestLoader(t *testing.T) {
	ch := &Chart{}

	f := ch.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func TestChart_render_without_title(t *testing.T) {
	ch := &Chart{
		logger: zap.NewNop(),
	}

	f := ch.render(nil)

	luaState := lua.NewState()
	n := f(luaState)

	assert.Equal(t, 2, n)
	e := luaState.Get(2)

	assert.Equal(t, "title must be defined", e.String())
}

func TestChart_render_without_params(t *testing.T) {
	ch := &Chart{
		logger: zap.NewNop(),
	}

	f := ch.render(nil)

	luaState := lua.NewState()
	luaState.Push(lua.LString("title"))
	n := f(luaState)

	assert.Equal(t, 2, n)
	e := luaState.Get(3)

	assert.Equal(t, "chart data table must be defined", e.String())
}

func TestChart_render_wrong_params(t *testing.T) {
	ch := &Chart{
		logger: zap.NewNop(),
	}

	f := ch.render(nil)

	luaState := lua.NewState()
	luaState.Push(lua.LString("title"))
	tt := &lua.LTable{}
	luaState.Push(tt)
	n := f(luaState)

	assert.Equal(t, 1, n)
}
