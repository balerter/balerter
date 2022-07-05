package slack

import (
	slackCfg "github.com/balerter/balerter/internal/config/channels/slack"
	"github.com/slack-go/slack"
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
	ignore  bool
}

// New creates new Slack channel
func New(cfg slackCfg.Slack, logger *zap.Logger) (*Slack, error) {
	m := &Slack{
		logger:  logger,
		name:    cfg.Name,
		channel: cfg.Channel,
		ignore:  cfg.Ignore,
	}

	m.api = slack.New(cfg.Token)

	return m, nil
}

// Name returns the channel name
func (m *Slack) Name() string {
	return m.name
}

func (m *Slack) Ignore() bool {
	return m.ignore
}
