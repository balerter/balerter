package mock

import (
	"crypto/md5"
	"fmt"
	"github.com/balerter/balerter/internal/lua_formatter"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type ModuleMock struct {
	name    string
	methods []string
	logger  *zap.Logger

	responses        map[string][]lua.LValue
	assertsCalled    map[string]*assert
	assertsNotCalled map[string]*assert
	queryLog         map[string]int
	errors           []string
}

type assert struct {
	method string
	args   []lua.LValue
	count  int
}

func New(name string, methods []string, logger *zap.Logger) *ModuleMock {
	m := &ModuleMock{
		name:             name,
		methods:          methods,
		logger:           logger,
		responses:        make(map[string][]lua.LValue),
		assertsCalled:    make(map[string]*assert),
		assertsNotCalled: make(map[string]*assert),
		queryLog:         make(map[string]int),
	}

	return m
}

func (m *ModuleMock) GetLoader(_ *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		exports := map[string]lua.LGFunction{
			"on":              m.on,
			"assertCalled":    m.assert(true),
			"assertNotCalled": m.assert(false),
		}

		for _, method := range m.methods {
			exports[method] = m.call(method)
		}

		mod := L.SetFuncs(L.NewTable(), exports)

		L.Push(mod)

		return 1
	}
}

func (m *ModuleMock) buildHash(methodName string, values []lua.LValue) string {
	s := methodName + ":"
	for _, v := range values {
		sv, err := lua_formatter.ValueToString(v)
		if err != nil {
			m.logger.Error("error marshal lua.Value to string", zap.Error(err))
			sv = "![ERROR:" + err.Error() + "]" // todo: return an error and doesn't build a hash?
		}
		s += v.Type().String() + ":" + sv
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func (m *ModuleMock) Name() string {
	return m.name
}

func (m *ModuleMock) Clean() {
	for key := range m.responses {
		delete(m.responses, key)
	}
	for key := range m.assertsCalled {
		delete(m.assertsCalled, key)
	}
	for key := range m.assertsNotCalled {
		delete(m.assertsNotCalled, key)
	}
	for key := range m.queryLog {
		delete(m.queryLog, key)
	}
}
