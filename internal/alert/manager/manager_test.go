package manager

import (
	"testing"

	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func TestManager_Init(t *testing.T) {
	m := New(nil, zap.NewNop())

	cfg := config.Channels{
		Email:    []*config.ChannelEmail{{Name: "email1"}},
		Slack:    []*config.ChannelSlack{{Name: "slack1"}},
		Telegram: []*config.ChannelTelegram{{Name: "tg1"}},
		Syslog:   []*config.ChannelSyslog{{Name: "sl1", Network: "udp", Address: "127.0.0.1:2000"}},
		Notify:   []*config.ChannelNotify{{Name: "n1"}},
		Discord:  []*config.ChannelDiscord{{Name: "discord1"}},
	}

	err := m.Init(&cfg)
	require.NoError(t, err)
	require.Equal(t, 6, len(m.channels))

	c, ok := m.channels["email1"]
	require.True(t, ok)
	assert.Equal(t, "email1", c.Name())

	c, ok = m.channels["slack1"]
	require.True(t, ok)
	assert.Equal(t, "slack1", c.Name())

	c, ok = m.channels["tg1"]
	require.True(t, ok)
	assert.Equal(t, "tg1", c.Name())

	c, ok = m.channels["sl1"]
	require.True(t, ok)
	assert.Equal(t, "sl1", c.Name())

	c, ok = m.channels["n1"]
	require.True(t, ok)
	assert.Equal(t, "n1", c.Name())

	c, ok = m.channels["discord1"]
	require.True(t, ok)
	assert.Equal(t, "discord1", c.Name())
}

func TestManager_Loader(t *testing.T) {
	m := New(nil, zap.NewNop())

	L := lua.NewState()

	f := m.GetLoader(&script.Script{})
	c := f(L)
	assert.Equal(t, 1, c)

	v := L.Get(1).(*lua.LTable)

	assert.IsType(t, &lua.LNilType{}, v.RawGet(lua.LString("wrong-name")))

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func TestManager_Name(t *testing.T) {
	m := &Manager{}

	assert.Equal(t, "alert", m.Name())
}

func TestManager_Stop(t *testing.T) {
	m := &Manager{}

	assert.NoError(t, m.Stop())
}
