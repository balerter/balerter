package email

import (
	"github.com/balerter/balerter/internal/config/channels/email"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSend(t *testing.T) {
	e := &Email{
		conf: email.Email{
			Name:   "foo",
			From:   "gopher@example.net",
			To:     "foo1@example.com;foo2@example.com",
			Cc:     "foo3@example.com;foo4@example.com",
			Host:   "localhost",
			Port:   "1025",
			Secure: "none",
		},

		name:   "email_test",
		logger: zap.NewNop(),
	}

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(e.conf.Host, e.conf.Port), timeout)
	if err != nil {
		t.Fatalf("error dial, %s", err)
		return
	}
	if conn != nil {
		err2 := conn.Close()
		assert.NoError(t, err2)
	}

	msg := &message.Message{
		Level:     "error",
		AlertName: "alert-id",
		Text:      "alert text",
		Image:     "alert image",
	}

	err = e.Send(msg)
	require.NoError(t, err)
}
