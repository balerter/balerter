package prometheus

import (
	"github.com/balerter/balerter/internal/config"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Prometheus struct {
	logger            *zap.Logger
	name              string
	url               string
	basicAuthUsername string
	basicAuthPassword string
	client            http.Client
}

func New(cfg config.DataSourcePrometheus, logger *zap.Logger) (*Prometheus, error) {
	m := &Prometheus{
		logger:            logger,
		name:              "prometheus." + cfg.Name,
		url:               cfg.URL,
		basicAuthUsername: cfg.BasicAuth.Username,
		basicAuthPassword: cfg.BasicAuth.Password,
	}
	m.client = http.Client{
		Timeout: time.Second * 30,
	}

	return m, nil
}

func (m *Prometheus) Stop() error {
	m.client.CloseIdleConnections()
	return nil
}

func (m *Prometheus) Name() string {
	return m.name
}

func (m *Prometheus) GetLoader() lua.LGFunction {
	return m.loader
}

func (m *Prometheus) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"querySingle": m.querySingle,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}
