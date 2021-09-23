package meta

import (
	"github.com/balerter/balerter/internal/modules"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// ModuleName returns the module name
func ModuleName() string {
	return "meta"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"priorExecutionTime",
	}
}

// Meta represents the Meta core module
type Meta struct {
	logger *zap.Logger
}

// New creates new Meta core module
func New(logger *zap.Logger) *Meta {
	m := &Meta{
		logger: logger,
	}

	return m
}

// Name returns the module name
func (m *Meta) Name() string {
	return ModuleName()
}

// GetLoader returns the lua loader
func (m *Meta) GetLoader(j modules.Job) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"priorExecutionTime": m.priorExecutionTime(j),
				"cronLocation":       m.getCronLocation(j),
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

// Stop the module
func (m *Meta) Stop() error {
	return nil
}

func (m *Meta) priorExecutionTime(j modules.Job) lua.LGFunction {
	return func(luaState *lua.LState) int {
		luaState.Push(lua.LNumber(j.GetPriorExecutionTime().Seconds()))
		return 1
	}
}

func (m *Meta) getCronLocation(j modules.Job) lua.LGFunction {
	return func(luaState *lua.LState) int {
		luaState.Push(lua.LString(j.GetCronLocation().String()))
		return 1
	}
}
