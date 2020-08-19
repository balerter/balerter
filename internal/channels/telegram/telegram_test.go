package telegram

import (
	config "github.com/balerter/balerter/internal/config/channels/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	TestToken               = "153667468:AAHlSHlMqSt1f_uFmVRJbm5gntu2HI4WW8I"
	ChatID                  = 76918703
)

func TestNew(t *testing.T) {
	tg, err := New(&config.Telegram{
		Name: "foo",
		ChatID: ChatID,
		Token: TestToken,
	}, nil)
	require.NoError(t, err)
	assert.IsType(t, &Telegram{}, tg)
	assert.Equal(t, "foo", tg.name)
	assert.Equal(t, int64(ChatID), tg.chatID)
}

func TestName(t *testing.T) {
	tg := &Telegram{name: "foo"}
	assert.Equal(t, "foo", tg.Name())
}
