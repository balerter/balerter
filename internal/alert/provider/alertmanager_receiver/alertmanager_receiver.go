package alertmanagerreceiver

import (
	"github.com/balerter/balerter/internal/alert/provider/webhook"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
)

type AMReceiver struct {
	name   string
	whCore *webhook.Core
	logger *zap.Logger
}

func New(cfg *config.ChannelAlertmanagerReceiver, logger *zap.Logger) (*AMReceiver, error) {
	cfg.Settings.Headers["content-type"] = "application/json"

	a := &AMReceiver{
		name:   cfg.Name,
		logger: logger,
		whCore: webhook.NewCore(cfg.Settings),
	}

	return a, nil
}

func (a *AMReceiver) Name() string {
	return a.name
}
