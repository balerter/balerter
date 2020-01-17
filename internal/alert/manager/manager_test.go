package manager

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestManager_Init(t *testing.T) {

	m := New(zap.NewNop())

	cfg := config.Channels{
		Slack: []config.ChannelSlack{
			{
				Name:                 "slack1",
				URL:                  "url",
				MessagePrefixSuccess: "success",
				MessagePrefixError:   "error",
			},
		},
	}

	err := m.Init(cfg)
	require.NoError(t, err)
	require.Equal(t, 1, len(m.channels))

	_, ok := m.channels["slack1"]
	require.True(t, ok)
}

func TestManager_Loader(t *testing.T) {
	m := New(zap.NewNop())

	L := lua.NewState()

	c := m.loader(L)
	assert.Equal(t, 1, c)

	v := L.Get(1).(*lua.LTable)

	assert.IsType(t, &lua.LNilType{}, v.RawGet(lua.LString("wrong-name")))

	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("on")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("off")))
}

func TestManager_getAlertName(t *testing.T) {
	m := New(zap.NewNop())
	var ok bool
	var name string
	var L *lua.LState

	L = lua.NewState()
	_, ok = m.getAlertName(L)
	require.False(t, ok)

	L = lua.NewState()
	L.Push(lua.LString("  "))
	_, ok = m.getAlertName(L)
	require.False(t, ok)

	L = lua.NewState()
	L.Push(lua.LString(" name "))
	name, ok = m.getAlertName(L)
	require.True(t, ok)
	assert.Equal(t, "name", name)
}
