package alertmanagerreceiver

import (
	"github.com/balerter/balerter/internal/config/channels/alertmanagerreceiver"
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	a, err := New(alertmanagerreceiver.AlertmanagerReceiver{Settings: webhook.Settings{Headers: map[string]string{}}}, "", nil)
	require.NoError(t, err)
	assert.IsType(t, &AMReceiver{}, a)
}

func TestName(t *testing.T) {
	a := &AMReceiver{name: "foo"}
	assert.Equal(t, "foo", a.Name())
}

func TestAMReceiver_Ignore(t *testing.T) {
	a := &AMReceiver{ignore: false}
	assert.False(t, a.Ignore())
}
