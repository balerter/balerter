package webhook

import (
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"net/http"
	"time"
)

const (
	defaultTimeout = 3000
)

// Core represent the channel of the type Webhook
type Core struct {
	cfg     webhook.Settings
	client  *http.Client
	timeout time.Duration
}

// NewCore creates new Core channel
func NewCore(cfg webhook.Settings) *Core {
	t := cfg.Timeout
	if t == 0 {
		t = defaultTimeout
	}

	timeout := time.Millisecond * time.Duration(t)

	client := &http.Client{
		Timeout: timeout,
	}

	c := &Core{
		cfg:     cfg,
		client:  client,
		timeout: timeout,
	}

	return c
}
