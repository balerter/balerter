package twiliovoice

import (
	"github.com/balerter/balerter/internal/config/channels/twiliovoice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	tw, err := New(twiliovoice.Twilio{}, zap.NewNop())
	require.NoError(t, err)
	assert.IsType(t, &TwilioVoice{}, tw)
}

func TestName(t *testing.T) {
	tw := &TwilioVoice{name: "tw"}
	assert.Equal(t, "tw", tw.Name())
}
