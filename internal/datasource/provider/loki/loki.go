package loki

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/balerter/balerter/internal/config/datasources/loki"
	"github.com/balerter/balerter/internal/modules"

	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

//go:generate moq -out module_mock_http_client.go -skip-ensure -fmt goimports . httpClient

var (
	defaultTimeout = time.Second * 5
)

// ModuleName returns the module name
func ModuleName(name string) string {
	return "loki." + name
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

// Loki represents the datasource of type Loki
type Loki struct {
	logger            *zap.Logger
	name              string
	url               *url.URL
	basicAuthUsername string
	basicAuthPassword string
	client            httpClient
	timeout           time.Duration
}

// New creates new Loki datasource
func New(cfg loki.Loki, logger *zap.Logger) (*Loki, error) {
	m := &Loki{
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
		Timeout: time.Second * 30,
	}

	return m, nil
}

func (m *Loki) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	return nil, http.StatusNotImplemented, fmt.Errorf("not implemented")
}

// Stop the datasource
func (m *Loki) Stop() error {
	m.client.CloseIdleConnections()
	return nil
}

// Name returns the datasource name
func (m *Loki) Name() string {
	return m.name
}

// GetLoader returns the datasource lua loader
func (m *Loki) GetLoader(_ modules.Job) lua.LGFunction {
	return m.loader
}

func (m *Loki) loader(luaState *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.doQuery,
		"range": m.doRange,
	}

	mod := luaState.SetFuncs(luaState.NewTable(), exports)

	luaState.Push(mod)
	return 1
}
