package webhook

import (
	"github.com/balerter/balerter/internal/config/channels/webhook"
	"net/http"
	"time"
)

const (
	defaultTimeout = 3000
)

type Core struct {
	cfg     webhook.Settings
	client  *http.Client
	timeout time.Duration
}

func NewCore(cfg webhook.Settings) *Core {
	t := cfg.Timeout
	if t == 0 {
		t = defaultTimeout
	}

	timeout := time.Millisecond * t

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
