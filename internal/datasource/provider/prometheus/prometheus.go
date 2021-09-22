package prometheus

import (
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"github.com/balerter/balerter/internal/modules"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

// ModuleName returns the module name
func ModuleName(name string) string {
	return "prometheus." + name
}

// Methods returns module methods
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

// Prometheus represents the datasource of the type Prometheus
type Prometheus struct {
	logger            *zap.Logger
	name              string
	url               *url.URL
	basicAuthUsername string
	basicAuthPassword string
	client            httpClient
	timeout           time.Duration
}

// New creates new Prometheus datasource
func New(cfg prometheus.Prometheus, logger *zap.Logger) (*Prometheus, error) {
	m := &Prometheus{
		logger:  logger,
		name:    ModuleName(cfg.Name),
		timeout: time.Millisecond * time.Duration(cfg.Timeout),
	}

	if cfg.BasicAuth != nil {
		m.basicAuthUsername = cfg.BasicAuth.Username
		m.basicAuthPassword = cfg.BasicAuth.Password
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

// Stop the datasource
func (m *Prometheus) Stop() error {
	m.client.CloseIdleConnections()
	return nil
}

// Name returns the datasource name
func (m *Prometheus) Name() string {
	return m.name
}

// GetLoader returns the datasource lua loader
func (m *Prometheus) GetLoader(_ modules.Job) lua.LGFunction {
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
