package webhook

import (
	"net/http"

	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
)

// Webhook implements a Provider for webhook notifications.
type Webhook struct {
	conf   *config.ChannelWebhook
	logger *zap.Logger
	client *http.Client
	name   string
}

func New(cfg config.ChannelWebhook, logger *zap.Logger) (*Webhook, error) {
	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	return &Webhook{
		conf:   &cfg,
		logger: logger,
		client: client,
		name:   cfg.Name,
	}, nil
}

func (w *Webhook) Name() string {
	return w.name
}
