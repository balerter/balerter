package slack

import (
	slackCfg "github.com/balerter/balerter/internal/config/channels/slack"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

// API is an interface of Slack API
type API interface {
	SendMessage(channel string, options ...slack.MsgOption) (string, string, string, error)
}

// Slack represents a channel of type Slack
type Slack struct {
	logger  *zap.Logger
	name    string
	channel string
	api     API
}

// New creates new Slack channel
func New(cfg slackCfg.Slack, logger *zap.Logger) (*Slack, error) {
	m := &Slack{
		logger:  logger,
		name:    cfg.Name,
		channel: cfg.Channel,
	}

	m.api = slack.New(cfg.Token)

	return m, nil
}

// Name returns the channel name
func (m *Slack) Name() string {
	return m.name
}
