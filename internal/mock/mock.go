package mock

import (
	"github.com/balerter/balerter/internal/mock/registry"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

const (
	AnyValue = "__TEST_ANY_VALUE__"
)

type Registry interface {
	Clean()
	Register(AnyValue, method string, callArgs, retArgs []lua.LValue) error
	Response(AnyValue, method string, args []lua.LValue) ([]lua.LValue, error)
	AddAssert(method string, args []lua.LValue, called bool) error
	AddCall(method string, args []lua.LValue) error
	Result() []modules.TestResult
}

type ModuleMock struct {
	name    string
	methods []string
	logger  *zap.Logger

	registry Registry

	errors []string
}

func New(name string, methods []string, logger *zap.Logger) *ModuleMock {
	m := &ModuleMock{
		name:     name,
		methods:  methods,
		logger:   logger,
		registry: registry.New(),
	}

	return m
}

func (m *ModuleMock) Stop() error {
	return nil
}

func (m *ModuleMock) GetLoader(_ *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		exports := map[string]lua.LGFunction{
			"on":              m.on,
			"assertCalled":    m.assert(true),
			"assertNotCalled": m.assert(false),
		}

		for _, method := range m.methods {
			exports[method] = m.call(method)
		}

		mod := luaState.SetFuncs(luaState.NewTable(), exports)

		luaState.Push(mod)

		return 1
	}
}

func (m *ModuleMock) Name() string {
	return m.name
}

func (m *ModuleMock) Clean() {
	m.registry.Clean()
	m.errors = m.errors[:0]
}
