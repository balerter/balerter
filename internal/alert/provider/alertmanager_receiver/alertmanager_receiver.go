package alertmanagerreceiver

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"net/url"
)

type AMReceiver struct {
	name   string
	url    string
	logger *zap.Logger
}

func New(cfg *config.ChannelAlertmanagerReceiver, logger *zap.Logger) (*AMReceiver, error) {
	a := &AMReceiver{
		name:   cfg.Name,
		logger: logger,
	}

	u, err := url.Parse(cfg.Settings.URL)
	if err != nil {
		return nil, fmt.Errorf("error parse url, %w", err)
	}

	a.url = u.String()

	return a, nil
}

func (a *AMReceiver) Name() string {
	return a.name
}
