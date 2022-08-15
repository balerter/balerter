package log

import (
	logConfig "github.com/balerter/balerter/internal/config/channels/log"
	"github.com/balerter/balerter/internal/message"

	"go.uber.org/zap"
)

// Log represents a channel of type Log
type Log struct {
	name   string
	logger *zap.Logger
	ignore bool
}

// New creates new log channel
func New(cfg logConfig.Log, logger *zap.Logger) (*Log, error) {
	sl := &Log{
		name:   cfg.Name,
		logger: logger,
		ignore: cfg.Ignore,
	}

	return sl, nil
}

// Name returns the channel name
func (lg *Log) Name() string {
	return lg.name
}

func (lg *Log) Ignore() bool {
	return lg.ignore
}

// Send the message to the channel
func (lg *Log) Send(mes *message.Message) error {
	lg.logger.Info("Log channel message",
		zap.String("channel name", lg.name),
		zap.String("alert id", mes.AlertName),
		zap.String("level", mes.Level),
		zap.String("message", mes.Text),
		zap.String("image", mes.Image),
		zap.Any("fields", mes.Fields),
	)

	return nil
}
