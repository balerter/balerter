package email

import (
	"os"

	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
)

// Email implements a Provider for email notifications.
type Email struct {
	conf     *config.ChannelEmail
	hostname string
	logger   *zap.Logger
	name     string
}

func New(cfg config.ChannelEmail, logger *zap.Logger) (*Email, error) {
	h, err := os.Hostname()
	// Use localhost if os.Hostname() fails
	if err != nil {
		h = "localhost.localdomain"
	}
	return &Email{conf: &cfg, hostname: h, logger: logger, name: cfg.Name}, nil
}

func (e *Email) Name() string {
	return e.name
}
