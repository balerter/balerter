package runtime

import (
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
	logLevel     string
	isDebug      bool
	isOnce       bool
	withScript   string
	configSource string
	logger       *zap.Logger
}

func New(logLevel string, isDebug, isOnce bool, withScript, configSource string, logger *zap.Logger) *Runtime {
	m := &Runtime{
		logLevel:     logLevel,
		isDebug:      isDebug,
		isOnce:       isOnce,
		withScript:   withScript,
		configSource: configSource,
		logger:       logger,
	}

	return m
}

func (m *Runtime) Name() string {
	return ModuleName()
}

func (m *Runtime) GetLoader(_ *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(L *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"logLevel":     m.returnString(m.logLevel),
				"isDebug":      m.returnBool(m.isDebug),
				"isOnce":       m.returnBool(m.isOnce),
				"withScript":   m.returnString(m.withScript),
				"configSource": m.returnString(m.configSource),
			}

			mod := L.SetFuncs(L.NewTable(), exports)

			L.Push(mod)
			return 1
		}
	}()
}

func (m *Runtime) Stop() error {
	return nil
}

func (m *Runtime) returnBool(v bool) lua.LGFunction {
	return func(L *lua.LState) int {
		L.Push(lua.LBool(v))
		return 1
	}
}

func (m *Runtime) returnString(v string) lua.LGFunction {
	return func(L *lua.LState) int {
		L.Push(lua.LString(v))
		return 1
	}
}
