package slack

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/config"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
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

func (m *Slack) Send(level alert.Level, message *message.Message, chartData *chartModule.Data) error {

	var imageURL string // todo

	opts := createSlackMessageOptions(message.AlertName, message.Text, imageURL, message.Fields...)

	_channel, _timestamp, _text, err := m.api.SendMessage(m.channel, opts...)

	m.logger.Debug("send slack message", zap.Int("level", int(level)), zap.String("channel", _channel), zap.String("timestamp", _timestamp), zap.String("text", _text), zap.Error(err))

	return err
}
