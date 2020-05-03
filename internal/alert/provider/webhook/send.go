package webhook

import (
	"context"
	"net/http"
	"strings"

	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
)

func (w *Webhook) Send(message *message.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), w.conf.Timeout)
	defer cancel()

	req, err := w.request(ctx, message)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: w.conf.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
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
	}
	return req, nil
}

func interpolate(s string, message *message.Message) string {
	if message == nil {
		return s
	}
	s = strings.ReplaceAll(s, "$level", message.Level)
	s = strings.ReplaceAll(s, "$alert_name", message.AlertName)
	s = strings.ReplaceAll(s, "$text", message.Text)
	s = strings.ReplaceAll(s, "$image", message.Image)

	return s
}
