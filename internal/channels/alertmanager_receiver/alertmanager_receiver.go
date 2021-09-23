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

// AMReceiver represents a channel of type AlertManagerReceiver
type AMReceiver struct {
	name   string
	whCore webHookCore
	logger *zap.Logger
	ignore bool
}

// New creates new AlertManagerReceiver channel
func New(cfg alertmanagerreceiver.AlertmanagerReceiver, version string, logger *zap.Logger) (*AMReceiver, error) {
	//cfg.Settings.Headers["content-type"] = "application/json" // todo(negasus): init headers map?

	a := &AMReceiver{
		name:   cfg.Name,
		logger: logger,
		whCore: webhook.NewCore(cfg.Settings, version),
		ignore: cfg.Ignore,
	}

	return a, nil
}

// Name returns name of the channel
func (a *AMReceiver) Name() string {
	return a.name
}

func (a *AMReceiver) Ignore() bool {
	return a.ignore
}
