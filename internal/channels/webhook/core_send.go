package webhook

import (
	"context"
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"github.com/balerter/balerter/internal/message"
	"io"
	"net/http"
)

// Send the message to the channel
func (w *Core) Send(body io.Reader, m *message.Message) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, w.cfg.Method, w.cfg.URL, body)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for param, value := range w.cfg.Payload.QueryParams {
		query.Add(param, interpolate(value, m))
	}
	req.URL.RawQuery = query.Encode()

	req.Header.Set("User-Agent", "Balerter "+w.version)

	for key, value := range w.cfg.Headers {
		req.Header.Add(key, value)
	}

	if w.cfg.Auth != nil {
		switch w.cfg.Auth.Type {
		case webhook.AuthTypeBasic:
			req.SetBasicAuth(w.cfg.Auth.AuthBasicConfig.Login, w.cfg.Auth.AuthBasicConfig.Password)
		case webhook.AuthTypeBearer:
			token := "Bearer " + w.cfg.Auth.AuthBearerConfig.Token
			req.Header.Add("Authorization", token)
		case webhook.AuthTypeCustom:
			for key, value := range w.cfg.Auth.AuthCustomConfig.Headers {
				req.Header.Add(key, value)
			}

			query := req.URL.Query()
			for param, value := range w.cfg.Auth.AuthCustomConfig.QueryParams {
				query.Add(param, value)
			}
			req.URL.RawQuery = query.Encode()
		}
	}

	return w.client.Do(req)
}
