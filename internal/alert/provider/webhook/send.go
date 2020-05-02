package webhook

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/balerter/balerter/internal/alert/message"
)

func (w *Webhook) Send(message *message.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := w.request(ctx, message)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (w *Webhook) request(ctx context.Context, message *message.Message) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)
	if w.conf.Method == "get" {
		req, err = w.getRequest(ctx, message)
	} else {
		req, err = w.postRequest(ctx, message)
	}
	if err != nil {
		return nil, err
	}

	switch w.conf.Auth.Type {
	case "basic":
		req.SetBasicAuth(w.conf.Auth.AuthBasicConfig.Login, w.conf.Auth.AuthBasicConfig.Password)
	case "bearer":
		token := "Bearer " + w.conf.Auth.AuthBearerConfig.Token
		req.Header.Add("Authorization", token)
	case "custom":
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

func (w *Webhook) getRequest(ctx context.Context, message *message.Message) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.conf.URL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for param, value := range w.conf.Payload.QueryParams {
		query.Add(param, w.interpolate(value, message))
	}

	return req, err
}

func (w *Webhook) postRequest(ctx context.Context, message *message.Message) (*http.Request, error) {
	body := w.interpolate(w.conf.Payload.Body, message)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.conf.URL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for param, value := range w.conf.Payload.QueryParams {
		query.Add(param, w.interpolate(value, message))
	}

	return req, err
}

func (w *Webhook) interpolate(s string, message *message.Message) string {
	if message == nil {
		return s
	}
	s = strings.ReplaceAll(s, "$level", message.Level)
	s = strings.ReplaceAll(s, "$alert_name", message.AlertName)
	s = strings.ReplaceAll(s, "$text", message.Text)
	s = strings.ReplaceAll(s, "$image", message.Image)

	return s
}
