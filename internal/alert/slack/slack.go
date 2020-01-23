package slack

import (
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

type Slack struct {
	logger               *zap.Logger
	name                 string
	channel              string
	messagePrefixSuccess string
	messagePrefixError   string
	api                  *slack.Client
}

func New(cfg config.ChannelSlack, logger *zap.Logger) (*Slack, error) {
	m := &Slack{
		logger:               logger,
		name:                 cfg.Name,
		channel:              cfg.Channel,
		messagePrefixSuccess: cfg.MessagePrefixSuccess,
		messagePrefixError:   cfg.MessagePrefixError,
	}

	m.api = slack.New(cfg.Token)

	return m, nil
}

func (m *Slack) Name() string {
	return m.name
}

func (m *Slack) SendSuccess(message *message.Message) error {
	msgOptions := createSlackMessageOptions(message.AlertName, m.messagePrefixSuccess+message.Text, message.Fields...)

	return m.send(msgOptions)
}

func (m *Slack) SendError(message *message.Message) error {
	msgOptions := createSlackMessageOptions(message.AlertName, m.messagePrefixError+message.Text, message.Fields...)

	return m.send(msgOptions)
}

func (m *Slack) Send(message *message.Message) error {
	msgOptions := createSlackMessageOptions(message.AlertName, message.Text, message.Fields...)

	return m.send(msgOptions)
}

func (m *Slack) send(opts []slack.MsgOption) error {
	_channel, _timestamp, _text, err := m.api.SendMessage(m.channel, opts...)

	m.logger.Debug("send slack message", zap.String("channel", _channel), zap.String("timestamp", _timestamp), zap.String("text", _text), zap.Error(err))

	return err
}
