package loki

import (
	"github.com/balerter/balerter/internal/config/datasources/loki"
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
	return "loki." + name
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

type Loki struct {
	logger            *zap.Logger
	name              string
	url               *url.URL
	basicAuthUsername string
	basicAuthPassword string
	client            httpClient
	timeout           time.Duration
}

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

func (m *Loki) Stop() error {
	m.client.CloseIdleConnections()
	return nil
}

func (m *Loki) Name() string {
	return m.name
}

func (m *Loki) GetLoader(_ *script.Script) lua.LGFunction {
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
