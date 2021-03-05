package alertmanager

import (
	"github.com/balerter/balerter/internal/channels/webhook"
	"github.com/balerter/balerter/internal/config/channels/alertmanager"
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type webHookCore interface {
	Send(body io.Reader, m *message.Message) (*http.Response, error)
}

type AlertManager struct {
	name   string
	whCore webHookCore
	logger *zap.Logger
}

func New(cfg alertmanager.Alertmanager, logger *zap.Logger) (*AlertManager, error) {
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
