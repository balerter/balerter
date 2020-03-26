package test

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type Mock struct {
	name   string
	logger *zap.Logger

	queryLog  map[string]int
	responses map[string][]lua.LValue
	errors    []string

	asserts struct {
		callQuery map[string]int
	}
}

func New(cfg config.DataSourceClickhouse, logger *zap.Logger) (*Mock, error) {
	m := &Mock{
		name:      "clickhouse." + cfg.Name,
		logger:    logger,
		queryLog:  make(map[string]int),
		responses: make(map[string][]lua.LValue),
		asserts: struct {
			callQuery map[string]int
		}{
			callQuery: make(map[string]int),
		},
	}

	return m, nil
}

func (m *Mock) Name() string {
	return m.name
}

func (m *Mock) GetLoader(_ *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"query":             m.query,
			"onQuery":           m.onQuery,
			"assertCalledQuery": m.assertCalledQuery,
		}

		mod := L.SetFuncs(L.NewTable(), exports)

		L.Push(mod)
		return 1
	}
}
