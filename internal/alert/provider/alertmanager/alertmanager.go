package alertmanager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"net/url"
)

const (
	VersionV1 = "v1"
	VersionV2 = "v2"

	defaultVersion = VersionV1
)

type AlertManager struct {
	name    string
	url     string
	version string
	logger  *zap.Logger
}

func New(cfg *config.ChannelAlertmanager, logger *zap.Logger) (*AlertManager, error) {
	a := &AlertManager{
		name:    cfg.Name,
		version: cfg.Version,
		logger:  logger,
	}

	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error parse url, %w", err)
	}

	if a.version == "" {
		a.version = defaultVersion
	}

	switch a.version {
	case VersionV1:
		u.Path = "/api/v1/alerts"
	case VersionV2:
		u.Path = "/api/v2/alerts"
	default:
		return nil, fmt.Errorf("unsuppored api version %s", a.version)
	}

	a.url = u.String()

	return a, nil
}

func (a *AlertManager) Name() string {
	return a.name
}
