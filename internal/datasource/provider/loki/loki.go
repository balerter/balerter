package loki

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
)

type Loki struct {
	logger            *zap.Logger
	name              string
	url               *url.URL
	basicAuthUsername string
	basicAuthPassword string
	client            http.Client
}

func New(cfg config.DataSourceLoki, logger *zap.Logger) (*Loki, error) {
	m := &Loki{
		logger:            logger,
		name:              "loki." + cfg.Name,
		basicAuthUsername: cfg.BasicAuth.Username,
		basicAuthPassword: cfg.BasicAuth.Password,
	}

	var err error

	m.url, err = url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	m.client = http.Client{
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

func (m *Loki) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.doQuery,
		"range": m.doRange,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}
