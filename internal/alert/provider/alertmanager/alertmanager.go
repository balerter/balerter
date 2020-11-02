package alertmanager

import (
	"github.com/balerter/balerter/internal/alert/provider/webhook"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
)

type AlertManager struct {
	name   string
	whCore *webhook.Core
	logger *zap.Logger
}

func New(cfg *config.ChannelAlertmanager, logger *zap.Logger) (*AlertManager, error) {
	a := &AlertManager{
		name:   cfg.Name,
		logger: logger,
		whCore: webhook.NewCore(cfg.Settings),
	}

	return a, nil
}

func (a *AlertManager) Name() string {
	return a.name
}
