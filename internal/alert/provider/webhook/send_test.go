package webhook

import (
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
)

func webhookSend(conf *config.ChannelWebhook, msg *message.Message) error {
	webhook, err := New(conf, zap.NewNop())
	if err != nil {
		return err
	}

	return webhook.Send(msg)
}

func testHook(conf *config.ChannelWebhook, msg *message.Message, h http.HandlerFunc) error {
	s := httptest.NewServer(h)
	defer s.Close()

	conf.Settings.URL = s.URL
	return webhookSend(conf, msg)
}

func getBasicAuthHeader(login, pass string) string {
	req, _ := http.NewRequest("GET", "", nil)
	req.SetBasicAuth(login, pass)
	return req.Header.Get("Authorization")
}

func TestSend(t *testing.T) {
	msg := &message.Message{
		Level:     "error",
		AlertName: "alert-id",
		Text:      "alert text",
		Image:     "alert image",
	}

	conf := &config.ChannelWebhook{
		Name: "foo",
		Settings: &config.WebhookSettings{
			Method: http.MethodPost,
			Auth: config.AuthConfig{
				Type: config.AuthTypeNone,
			},
			Payload: config.PayloadConfig{
				Body: `{"message": "$text"}`,
			},
			Timeout: 5 * time.Second,
		},
	}

	t.Run("no-auth", func(t *testing.T) {
		a := require.New(t)

		err := testHook(conf, msg, func(w http.ResponseWriter, req *http.Request) {
			a.Equal(conf.Settings.Method, req.Method)

			b, err := ioutil.ReadAll(req.Body)
			a.NoError(err)
			a.Equal(`{"message": "alert text"}`, string(b))
		})
		a.NoError(err)
	})

	t.Run("query-params-payload", func(t *testing.T) {
		a := require.New(t)

		conf.Settings.Payload = config.PayloadConfig{
			QueryParams: map[string]string{
				"foo": "bar",
			},
		}

		err := testHook(conf, msg, func(w http.ResponseWriter, req *http.Request) {
			err := req.ParseForm()
			a.NoError(err)
			a.Equal("bar", req.Form.Get("foo"))
		})
		a.NoError(err)
	})

	t.Run("basic-auth", func(t *testing.T) {
		a := require.New(t)

		login, pass := "login", "pass"
		conf.Settings.Auth = config.AuthConfig{
			Type: config.AuthTypeBasic,
			AuthBasicConfig: config.AuthBasicConfig{
				Login:    login,
				Password: pass,
			},
		}

		err := testHook(conf, msg, func(w http.ResponseWriter, req *http.Request) {
			h := req.Header.Get("Authorization")
			a.Equal(getBasicAuthHeader(login, pass), h)
		})
		a.NoError(err)
	})

	t.Run("bearer", func(t *testing.T) {
		a := require.New(t)

		conf.Settings.Auth = config.AuthConfig{
			Type: config.AuthTypeBearer,
			AuthBearerConfig: config.AuthBearerConfig{
				Token: "test-token",
			},
		}

		err := testHook(conf, msg, func(w http.ResponseWriter, req *http.Request) {
			h := req.Header.Get("Authorization")
			a.Equal("Bearer test-token", h)
		})
		a.NoError(err)
	})

	t.Run("custom", func(t *testing.T) {
		a := require.New(t)
		conf.Settings.Auth = config.AuthConfig{
			Type: config.AuthTypeCustom,
			AuthCustomConfig: config.AuthCustomConfig{
				Headers: map[string]string{
					"X-Header": "foo",
				},
				QueryParams: map[string]string{
					"q": "bar",
				},
			},
		}

		err := testHook(conf, msg, func(w http.ResponseWriter, req *http.Request) {
			h := req.Header.Get("X-Header")
			a.Equal("foo", h)

			param := req.URL.Query().Get("q")
			a.Equal("bar", param)
		})
		a.NoError(err)
	})

	t.Run("timeout-error", func(t *testing.T) {
		a := require.New(t)
		conf.Settings.Timeout = 1000

		err := testHook(conf, msg, func(w http.ResponseWriter, req *http.Request) {
			time.Sleep(3 * time.Second)
		})
		a.Error(err)
	})
}

func Test_interpolate(t *testing.T) {
	format := "$level:$alert_name:$text:$image"

	t.Run("nil", func(t *testing.T) {
		r := interpolate(format, nil)
		require.Equal(t, format, r)
	})
	t.Run("non-nil", func(t *testing.T) {
		r := interpolate(format, &message.Message{
			Level:     "level",
			AlertName: "alert_name",
			Text:      "text",
			Image:     "image",
		})
		require.Equal(t, "level:alert_name:text:image", r)
	})
}
