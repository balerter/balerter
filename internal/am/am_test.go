package manager

import (
	"github.com/balerter/balerter/internal/config/channels"
	"github.com/balerter/balerter/internal/config/channels/discord"
	"github.com/balerter/balerter/internal/config/channels/email"
	"github.com/balerter/balerter/internal/config/channels/notify"
	"github.com/balerter/balerter/internal/config/channels/slack"
	"github.com/balerter/balerter/internal/config/channels/syslog"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestManager_Init(t *testing.T) {
	m := New(nil, zap.NewNop())

	cfg := channels.Channels{
		Email:    []*email.Email{{Name: "email1"}},
		Slack:    []*slack.Slack{{Name: "slack1"}},
		Telegram: []*telegram.Telegram{{Name: "tg1", ChatID: 76918703, Token: "153667468:AAHlSHlMqSt1f_uFmVRJbm5gntu2HI4WW8I"}},
		Syslog:   []*syslog.Syslog{{Name: "sl1", Network: "udp", Address: "127.0.0.1:2000"}},
		Notify:   []*notify.Notify{{Name: "n1"}},
		Discord:  []*discord.Discord{{Name: "discord1"}},
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
