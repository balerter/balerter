package slack

import (
	slackCfg "github.com/balerter/balerter/internal/config/channels/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	s, err := New(slackCfg.Slack{}, zap.NewNop())
	require.NoError(t, err)
	assert.IsType(t, &Slack{}, s)
}

func TestName(t *testing.T) {
	s := &Slack{name: "foo"}
	assert.Equal(t, "foo", s.Name())
}

func TestSlack_Ignore(t *testing.T) {
	s := &Slack{ignore: true}
	assert.True(t, s.Ignore())
}
