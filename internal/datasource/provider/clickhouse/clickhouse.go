package clickhouse

import (
	"github.com/balerter/balerter/internal/config"
	lua "github.com/yuin/gopher-lua"
)

type Clickhouse struct {
	name string
}

func New(cfg config.DataSourceClickhouse) (*Clickhouse, error) {
	c := &Clickhouse{
		name: cfg.Name,
	}

	return c, nil
}

func (m *Clickhouse) Name() string {
	return m.name
}

func (m *Clickhouse) GetLoader() lua.LGFunction {
	return m.loader
}

func (m *Clickhouse) loader(L *lua.LState) int {
	return 0
}
