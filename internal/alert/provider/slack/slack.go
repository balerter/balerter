package slack

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

type Slack struct {
	logger  *zap.Logger
	name    string
	channel string
	api     *slack.Client
}

func New(cfg config.ChannelSlack, logger *zap.Logger) (*Slack, error) {
	m := &Slack{
		logger:  logger,
		name:    cfg.Name,
		channel: cfg.Channel,
	}

	m.api = slack.New(cfg.Token)

	return m, nil
}

func (m *Slack) Name() string {
	return m.name
}
