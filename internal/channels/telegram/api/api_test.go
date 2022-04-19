package api

import (
	"testing"

	"github.com/balerter/balerter/internal/config/channels/telegram"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	a, err := New(telegram.Telegram{Token: "A"})
	require.NoError(t, err)
	assert.IsType(t, &API{}, a)
	assert.Equal(t, "https://api.telegram.org/botA/", a.endpoint)
}
