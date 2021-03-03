package runtime

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func ModuleName() string {
	return "runtime"
}

func Methods() []string {
	return []string{
		"logLevel",
		"isDebug",
		"isOnce",
		"withScript",
		"configSource",
	}
}

type Runtime struct {
	flg    *config.Flags
	logger *zap.Logger
}

func New(flg *config.Flags, logger *zap.Logger) *Runtime {
	m := &Runtime{
		flg:    flg,
		logger: logger,
	}

	return m
}

func (m *Runtime) Name() string {
	return ModuleName()
}

func (m *Runtime) GetLoader(_ *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"logLevel":     m.returnString(m.flg.LogLevel),
				"isDebug":      m.returnBool(m.flg.Debug),
				"isOnce":       m.returnBool(m.flg.Once),
				"withScript":   m.returnString(m.flg.Script),
				"configSource": m.returnString(m.flg.ConfigFilePath),
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

func (m *Runtime) Stop() error {
	return nil
}

func (m *Runtime) returnBool(v bool) lua.LGFunction {
	return func(luaState *lua.LState) int {
		luaState.Push(lua.LBool(v))
		return 1
	}
}

func (m *Runtime) returnString(v string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		luaState.Push(lua.LString(v))
		return 1
	}
}
