package discord

import (
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
)

// Discord implements a Provider for discord notifications.
type Discord struct {
	conf   *config.ChannelDiscord
	logger *zap.Logger
	name   string
}

// New returns the new Discord instance
func New(cfg *config.ChannelDiscord, logger *zap.Logger) (*Discord, error) {
	return &Discord{conf: cfg, logger: logger, name: cfg.Name}, nil
}

// Name returns the Discord channel name
func (d *Discord) Name() string {
	return d.name
}
