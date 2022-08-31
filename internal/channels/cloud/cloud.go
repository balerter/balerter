package cloud

import (
	"github.com/balerter/balerter/internal/cloud"
	cloudCfg "github.com/balerter/balerter/internal/config/channels/cloud"
	"github.com/balerter/balerter/internal/message"

	"go.uber.org/zap"
)

type Cloud struct {
	name   string
	logger *zap.Logger
	ignore bool
}

// New creates new log channel
func New(cfgCh cloudCfg.Cloud, logger *zap.Logger) (*Cloud, error) {
	sl := &Cloud{
		name:   cfgCh.Name,
		ignore: cfgCh.Ignore,
		logger: logger,
	}

	return sl, nil
}

// Name returns the channel name
func (cl *Cloud) Name() string {
	return cl.name
}

func (cl *Cloud) Ignore() bool {
	return cl.ignore
}

// Send the message to the channel
func (cl *Cloud) Send(mes *message.Message) error {
	return cloud.SendMessage(mes)
}
