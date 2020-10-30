package webhook

import (
	"net/http"
	"time"

	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
)

const (
	defaultTimeout = 3000
)

// Webhook implements a Provider for webhook notifications.
type Webhook struct {
	conf    *config.ChannelWebhook
	logger  *zap.Logger
	client  *http.Client
	name    string
	timeout time.Duration
}

func New(cfg *config.ChannelWebhook, logger *zap.Logger) (*Webhook, error) {
	t := cfg.Timeout
	if t == 0 {
		t = defaultTimeout
	}

	timeout := time.Millisecond * time.Duration(t)

	client := &http.Client{
		Timeout: timeout,
	}

	return &Webhook{
		conf:    cfg,
		logger:  logger,
		client:  client,
		name:    cfg.Name,
		timeout: timeout,
	}, nil
}

func (w *Webhook) Name() string {
	return w.name
}
