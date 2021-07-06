package manager

import (
	"github.com/balerter/balerter/internal/config/channels"
	"github.com/balerter/balerter/internal/config/channels/alertmanager"
	"github.com/balerter/balerter/internal/config/channels/alertmanagerreceiver"
	"github.com/balerter/balerter/internal/config/channels/discord"
	"github.com/balerter/balerter/internal/config/channels/email"
	"github.com/balerter/balerter/internal/config/channels/notify"
	"github.com/balerter/balerter/internal/config/channels/slack"
	"github.com/balerter/balerter/internal/config/channels/syslog"
	"github.com/balerter/balerter/internal/config/channels/telegram"
	"github.com/balerter/balerter/internal/config/channels/twiliovoice"
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestManager_Init(t *testing.T) {
	m := New(zap.NewNop())

	cfg := &channels.Channels{
		Email:                []email.Email{{Name: "email1"}},
		Slack:                []slack.Slack{{Name: "slack1"}},
		Telegram:             []telegram.Telegram{{Name: "tg1"}},
		Syslog:               []syslog.Syslog{{Name: "sl1", Network: "udp", Address: "127.0.0.1:2000"}},
		Notify:               []notify.Notify{{Name: "n1"}},
		Discord:              []discord.Discord{{Name: "discord1"}},
		Webhook:              []webhook.Webhook{{Name: "wh1"}},
		Alertmanager:         []alertmanager.Alertmanager{{Name: "am1"}},
		AlertmanagerReceiver: []alertmanagerreceiver.AlertmanagerReceiver{{Name: "amr1"}},
		TwilioVoice:          []twiliovoice.Twilio{{Name: "tw1"}},
	}

	err := m.Init(cfg)
	require.NoError(t, err)
	require.Equal(t, 10, len(m.channels))

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

	c, ok = m.channels["wh1"]
	require.True(t, ok)
	assert.Equal(t, "wh1", c.Name())

	c, ok = m.channels["am1"]
	require.True(t, ok)
	assert.Equal(t, "am1", c.Name())

	c, ok = m.channels["amr1"]
	require.True(t, ok)
	assert.Equal(t, "amr1", c.Name())

	c, ok = m.channels["tw1"]
	require.True(t, ok)
	assert.Equal(t, "tw1", c.Name())
}
