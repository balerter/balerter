package slack

import (
	"github.com/balerter/balerter/internal/message"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

type mockAPI struct {
	mock.Mock
}

func (m *mockAPI) SendMessage(channel string, options ...slack.MsgOption) (ch, timestamp, text string, err error) {
	args := m.Called(channel, options)

	ch = args.String(0)
	timestamp = args.String(1)
	text = args.String(2)
	err = args.Error(3)
	return
}

func TestSend(t *testing.T) {
	api := &mockAPI{}
	api.On("SendMessage", mock.Anything, mock.Anything).Return("1", "2", "3", nil)

	m := &Slack{
		api:    api,
		logger: zap.NewNop(),
	}

	mes := &message.Message{
		Level:     "error",
		AlertName: "alert-id",
		Text:      "alert text",
		Image:     "alert image",
		Fields:    map[string]string{"a": "b"},
	}

	err := m.Send(mes)
	require.NoError(t, err)
}
