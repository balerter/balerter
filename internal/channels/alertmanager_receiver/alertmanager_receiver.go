package alertmanagerreceiver

import (
	"github.com/balerter/balerter/internal/channels/webhook"
	"github.com/balerter/balerter/internal/config/channels/alertmanagerreceiver"
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type webHookCore interface {
	Send(body io.Reader, m *message.Message) (*http.Response, error)
}

type AMReceiver struct {
	name   string
	whCore webHookCore
	logger *zap.Logger
}

func New(cfg alertmanagerreceiver.AlertmanagerReceiver, logger *zap.Logger) (*AMReceiver, error) {
	//cfg.Settings.Headers["content-type"] = "application/json" // todo(negasus): init headers map?

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
