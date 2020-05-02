package webhook

import (
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
)

// Webhook implements a Provider for webhook notifications.
type Webhook struct {
	conf   *config.ChannelWebhook
	logger *zap.Logger
	name   string
}

func New(cfg config.ChannelWebhook, logger *zap.Logger) (*Webhook, error) {
	return &Webhook{conf: &cfg, logger: logger, name: cfg.Name}, nil
}

func (w *Webhook) Name() string {
	return w.name
}
