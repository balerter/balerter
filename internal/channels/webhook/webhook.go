package webhook

import (
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"go.uber.org/zap"
)

// Webhook implements a Provider for webhook notifications.
type Webhook struct {
	logger *zap.Logger
	name   string
	body   string
	whCore *Core
	ignore bool
}

// New creates new Webhook channel
func New(cfg webhook.Webhook, logger *zap.Logger) (*Webhook, error) {
	return &Webhook{
		body:   cfg.Settings.Payload.Body,
		logger: logger,
		name:   cfg.Name,
		whCore: NewCore(cfg.Settings),
		ignore: cfg.Ignore,
	}, nil
}

// Name returns the channel name
func (w *Webhook) Name() string {
	return w.name
}

func (w *Webhook) Ignore() bool {
	return w.ignore
}
