package webhook

import (
	"context"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSend(t *testing.T) {
	w := &Webhook{
		conf: &config.ChannelWebhook{
			Name: "foo",
			Auth: config.AuthConfig{
				Type: "basic",
				AuthBasicConfig: config.AuthBasicConfig{
					Login:    "login",
					Password: "password",
				},
			},
			Payload: config.PayloadConfig{
				Body: `{"message": "$text"}`,
			},
		},
		name:   "webhook_test",
		logger: zap.NewNop(),
	}

	msg := &message.Message{
		Level:     "error",
		AlertName: "alert-id",
		Text:      "alert text",
		Image:     "alert image",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := w.request(ctx, msg)
	switch w.conf.Auth.Type {
	case "basic":
		username, password, ok := req.BasicAuth()
		if !ok {
			t.Errorf("basic auth must be not empty")
		}
		if username != w.conf.Auth.AuthBasicConfig.Login {
			t.Errorf("username want: %s, got: %s", w.conf.Auth.AuthBasicConfig.Login, username)
		}
		if password != w.conf.Auth.AuthBasicConfig.Password {
			t.Errorf("password want: %s, got: %s", w.conf.Auth.AuthBasicConfig.Password, password)
		}
	case "bearer":
		want := "bearer " + w.conf.Auth.AuthBearerConfig.Token
		got := req.Header.Get("Authorization")
		if got != want {
			t.Errorf("bearer token want: %s, got: %s", want, got)
		}
	case "custom":
		for key, value := range w.conf.Auth.AuthCustomConfig.Headers {
			want := value
			got := req.Header.Get(key)
			if got != want {
				t.Errorf("custom header want: %s, got: %s", want, got)
			}
		}
		query := req.URL.Query()
		for param, value := range w.conf.Auth.AuthCustomConfig.QueryParams {
			want := value
			got := query.Get(param)
			if got != want {
				t.Errorf("custom query param want: %s, got: %s", want, got)
			}
		}
	}

	require.NoError(t, err)
}
