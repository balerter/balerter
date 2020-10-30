package webhook

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"net/http"
	"strings"
)

func (w *Webhook) Send(m *message.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	req, err := w.request(ctx, m)
	if err != nil {
		return fmt.Errorf("webhook request creation failed: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (w *Webhook) request(ctx context.Context, m *message.Message) (*http.Request, error) {
	body := interpolate(w.conf.Payload.Body, m)

	req, err := http.NewRequestWithContext(ctx, w.conf.Method, w.conf.URL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for param, value := range w.conf.Payload.QueryParams {
		query.Add(param, interpolate(value, m))
	}
	req.URL.RawQuery = query.Encode()

	switch w.conf.Auth.Type {
	case config.AuthTypeBasic:
		req.SetBasicAuth(w.conf.Auth.AuthBasicConfig.Login, w.conf.Auth.AuthBasicConfig.Password)
	case config.AuthTypeBearer:
		token := "Bearer " + w.conf.Auth.AuthBearerConfig.Token
		req.Header.Add("Authorization", token)
	case config.AuthTypeCustom:
		for key, value := range w.conf.Auth.AuthCustomConfig.Headers {
			req.Header.Add(key, value)
		}

		query := req.URL.Query()
		for param, value := range w.conf.Auth.AuthCustomConfig.QueryParams {
			query.Add(param, value)
		}
		req.URL.RawQuery = query.Encode()
	}
	return req, nil
}

func interpolate(s string, m *message.Message) string {
	if m == nil {
		return s
	}

	return strings.NewReplacer(
		"$level", m.Level,
		"$alert_name", m.AlertName,
		"$text", m.Text,
		"$image", m.Image,
	).Replace(s)
}
