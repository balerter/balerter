package telegram

import (
	"testing"

	"github.com/balerter/balerter/internal/channels/telegram/api"
	"github.com/balerter/balerter/internal/message"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSend_WithoutImage(t *testing.T) {
	var tgMessage *api.TextMessage

	m := &APIerMock{
		SendTextMessageFunc: func(textMessage *api.TextMessage) error {
			tgMessage = textMessage
			return nil
		},
	}

	tg := &Telegram{
		api:    m,
		logger: zap.NewNop(),
		chatID: 42,
	}

	mes := &message.Message{
		Level:     "foo",
		AlertName: "bar",
		Text:      "baz",
		Image:     "",
	}

	err := tg.Send(mes)
	require.NoError(t, err)

	assert.Equal(t, 1, len(m.SendTextMessageCalls()))
	require.NotNil(t, tgMessage)
	assert.Equal(t, "baz", tgMessage.Text)
	assert.Equal(t, int64(42), tgMessage.ChatID)
}

func TestSend_WithImage(t *testing.T) {
	var tgMessage *api.PhotoMessage

	m := &APIerMock{
		SendPhotoMessageFunc: func(photoMessage *api.PhotoMessage) error {
			tgMessage = photoMessage
			return nil
		},
	}
	tg := &Telegram{
		api:    m,
		logger: zap.NewNop(),
		chatID: 42,
	}

	mes := &message.Message{
		Level:     "foo",
		AlertName: "bar",
		Text:      "baz",
		Image:     "img1",
	}

	err := tg.Send(mes)
	require.NoError(t, err)

	assert.Equal(t, 1, len(m.SendPhotoMessageCalls()))

	require.NotNil(t, tgMessage)
	assert.Equal(t, int64(42), tgMessage.ChatID)
	assert.Equal(t, "baz", tgMessage.Caption)
	assert.Equal(t, "img1", tgMessage.Photo)
}
