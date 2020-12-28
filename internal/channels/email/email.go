package email

import (
	"github.com/balerter/balerter/internal/config/channels/email"
	"os"

	"go.uber.org/zap"
)

// Email implements a Provider for email notifications.
type Email struct {
	conf     *email.Email
	hostname string
	logger   *zap.Logger
	name     string
}

// New returns the new Email instance
func New(cfg *email.Email, logger *zap.Logger) (*Email, error) {
	h, err := os.Hostname()
	// Use localhost if os.Hostname() fails
	if err != nil {
		h = "localhost.localdomain"
	}
	return &Email{conf: cfg, hostname: h, logger: logger, name: cfg.Name}, nil
}

// Name returns the Email channel name
func (e *Email) Name() string {
	return e.name
}
