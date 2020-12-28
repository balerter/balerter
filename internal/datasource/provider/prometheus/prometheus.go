package prometheus

import (
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

func ModuleName(name string) string {
	return "prometheus." + name
}

func Methods() []string {
	return []string{
		"query",
		"range",
	}
}

type httpClient interface {
	CloseIdleConnections()
	Do(r *http.Request) (*http.Response, error)
}

type Prometheus struct {
	logger            *zap.Logger
	name              string
	url               *url.URL
	basicAuthUsername string
	basicAuthPassword string
	client            httpClient
	timeout           time.Duration
}

func New(cfg *prometheus.Prometheus, logger *zap.Logger) (*Prometheus, error) {
	m := &Prometheus{
		logger:            logger,
		name:              ModuleName(cfg.Name),
		basicAuthUsername: cfg.BasicAuth.Username,
		basicAuthPassword: cfg.BasicAuth.Password,
		timeout:           cfg.Timeout,
	}

	if m.timeout == 0 {
		m.timeout = defaultTimeout
	}

	var err error

	m.url, err = url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	m.client = &http.Client{
		Timeout: m.timeout,
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

func (m *Prometheus) GetLoader(_ *script.Script) lua.LGFunction {
	return m.loader
}

func (m *Prometheus) loader(luaState *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.doQuery,
		"range": m.doRange,
	}

	mod := luaState.SetFuncs(luaState.NewTable(), exports)

	luaState.Push(mod)
	return 1
}
