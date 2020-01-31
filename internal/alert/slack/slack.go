package slack

import (
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

type Slack struct {
	logger   *zap.Logger
	name     string
	channel  string
	api      *slack.Client
	prefixes config.ChannelPrefixes
}

func New(cfg config.ChannelSlack, logger *zap.Logger) (*Slack, error) {
	m := &Slack{
		logger:   logger,
		name:     cfg.Name,
		channel:  cfg.Channel,
		prefixes: cfg.Prefixes,
	}

	m.api = slack.New(cfg.Token)

	return m, nil
}

func (m *Slack) Name() string {
	return m.name
}

func (m *Slack) Send(level message.Level, message *message.Message) error {
	opts := createSlackMessageOptions(message.AlertName, message.Text, message.Fields...)

	_channel, _timestamp, _text, err := m.api.SendMessage(m.channel, opts...)

	m.logger.Debug("send slack message", zap.Int("level", int(level)), zap.String("channel", _channel), zap.String("timestamp", _timestamp), zap.String("text", _text), zap.Error(err))

	return err
}
