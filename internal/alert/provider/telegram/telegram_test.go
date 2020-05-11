package telegram

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	tg, err := New(&config.ChannelTelegram{Name: "foo", ChatID: 42}, nil)
	require.NoError(t, err)
	assert.IsType(t, &Telegram{}, tg)
	assert.Equal(t, "foo", tg.name)
	assert.Equal(t, int64(42), tg.chatID)
}

func TestName(t *testing.T) {
	tg := &Telegram{name: "foo"}
	assert.Equal(t, "foo", tg.Name())
}
