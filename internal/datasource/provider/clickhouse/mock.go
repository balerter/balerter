package clickhouse

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type Mock struct {
	name   string
	logger *zap.Logger

	responses map[string][]lua.LValue
	errors    []error
}

func NewMock(cfg config.DataSourceClickhouse, logger *zap.Logger) (*Mock, error) {
	m := &Mock{
		name:      "clickhouse." + cfg.Name,
		logger:    logger,
		responses: make(map[string][]lua.LValue),
	}

	return m, nil
}

func (m *Mock) Errors() []error {
	return m.errors
}

func (m *Mock) Stop() error {
	return nil
}

func (m *Mock) Name() string {
	return m.name
}

func (m *Mock) GetLoader(_ *script.Script) lua.LGFunction {
	return m.loader
}

func (m *Mock) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query":   m.query,
		"onQuery": m.onQuery,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}

func (m *Mock) writeResponse(query string) lua.LGFunction {
	return func(L *lua.LState) int {
		args := make([]lua.LValue, L.GetTop())
		for i := 0; i < L.GetTop(); i++ {
			args[i] = L.Get(i + 1) // lua indexing starts with 1
		}

		m.responses[query] = args

		return 0
	}
}

func (m *Mock) query(L *lua.LState) int {

	m.logger.Debug("call clickhouse mock query")

	if L.GetTop() != 1 {
		m.logger.Error("wrong arguments count")
		return 0
	}

	q := L.Get(1).String()

	args, ok := m.responses[q]
	if !ok {
		err := fmt.Errorf("response is not defined for the query: " + q)
		m.errors = append(m.errors, err)
		//m.logger.Error(err.Error())
		return 0
	}

	for _, a := range args {
		L.Push(a)
	}

	return len(args)
}

func (m *Mock) onQuery(L *lua.LState) int {

	if L.GetTop() == 0 {
		m.logger.Error("arguments not found")
		return 0
	}

	q := L.Get(1).String()

	//m.logger.Debug("clickhouse onQuery", zap.String("query", q), zap.String("datasource name", m.name))

	T := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{"response": m.writeResponse(q)})

	L.Push(T)

	return 1
}
