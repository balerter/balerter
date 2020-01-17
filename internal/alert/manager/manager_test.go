package manager

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/require"
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
