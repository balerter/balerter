package email

import (
	"testing"

	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSend(t *testing.T) {

	e := &Email{
		conf: &config.ChannelEmail{
			Name:       "foo",
			From:       "gopher@example.net",
			To:         "foo@example.com",
			ServerName: "localhost",
			ServerPort: "1025",
		},

		name:   "email_test",
		logger: zap.NewNop(),
	}

	msg := &message.Message{
		Level:     "error",
		AlertName: "alert-id",
		Text:      "alert text",
		Fields:    []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11"},
		Image:     "alert image",
	}

	err := e.Send(msg)
	require.NoError(t, err)
}
