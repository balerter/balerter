package telegram

import (
	"github.com/balerter/balerter/internal/channels/telegram/api"
	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

type apiMock struct {
	mock.Mock
}

func (m *apiMock) SendTextMessage(mes *api.TextMessage) error {
	args := m.Called(mes)
	return args.Error(0)
}

func (m *apiMock) SendPhotoMessage(mes *api.PhotoMessage) error {
	args := m.Called(mes)
	return args.Error(0)
}

func TestSend_WithoutImage(t *testing.T) {
	var tgMessage *api.TextMessage

	m := &apiMock{}
	m.On("SendTextMessage", mock.Anything).Run(func(args mock.Arguments) {
		tgMessage = args.Get(0).(*api.TextMessage)
	}).Return(nil)

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

	m.AssertCalled(t, "SendTextMessage", mock.Anything)
	require.NotNil(t, tgMessage)
	assert.Equal(t, "baz", tgMessage.Text)
	assert.Equal(t, int64(42), tgMessage.ChatID)
}

func TestSend_WithImage(t *testing.T) {
	var tgMessage *api.PhotoMessage

	m := &apiMock{}
	m.On("SendTextMessage", mock.Anything).Return(nil)
	m.On("SendPhotoMessage", mock.Anything).Run(func(args mock.Arguments) {
		tgMessage = args.Get(0).(*api.PhotoMessage)
	}).Return(nil)

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

	m.AssertCalled(t, "SendTextMessage", mock.Anything)
	m.AssertCalled(t, "SendPhotoMessage", mock.Anything)
	require.NotNil(t, tgMessage)
	assert.Equal(t, int64(42), tgMessage.ChatID)
	assert.Equal(t, "", tgMessage.Caption)
	assert.Equal(t, "img1", tgMessage.Photo)
}
