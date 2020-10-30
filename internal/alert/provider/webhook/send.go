package webhook

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"net/http"
	"strings"
)

func (w *Webhook) Send(message *message.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	req, err := w.request(ctx, message)
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

func (w *Webhook) request(ctx context.Context, message *message.Message) (*http.Request, error) {
	body := interpolate(w.conf.Payload.Body, message)

	req, err := http.NewRequestWithContext(ctx, w.conf.Method, w.conf.URL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for param, value := range w.conf.Payload.QueryParams {
		query.Add(param, interpolate(value, message))
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

func interpolate(s string, message *message.Message) string {
	if message == nil {
		return s
	}

	return strings.NewReplacer(
		"$level", message.Level,
		"$alert_name", message.AlertName,
		"$text", message.Text,
		"$image", message.Image,
	).Replace(s)
}
