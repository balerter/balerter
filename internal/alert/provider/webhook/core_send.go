package webhook

import (
	"context"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"io"
	"net/http"
)

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

	for key, value := range w.cfg.Headers {
		req.Header.Add(key, value)
	}

	switch w.cfg.Auth.Type {
	case config.AuthTypeBasic:
		req.SetBasicAuth(w.cfg.Auth.AuthBasicConfig.Login, w.cfg.Auth.AuthBasicConfig.Password)
	case config.AuthTypeBearer:
		token := "Bearer " + w.cfg.Auth.AuthBearerConfig.Token
		req.Header.Add("Authorization", token)
	case config.AuthTypeCustom:
		for key, value := range w.cfg.Auth.AuthCustomConfig.Headers {
			req.Header.Add(key, value)
		}

		query := req.URL.Query()
		for param, value := range w.cfg.Auth.AuthCustomConfig.QueryParams {
			query.Add(param, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	return w.client.Do(req)
}
