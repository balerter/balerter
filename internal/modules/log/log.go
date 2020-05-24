package log

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type Log struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Log {
	l := &Log{
		logger: logger,
	}

	return l
}

func (l *Log) Name() string {
	return ModuleName()
}

func (l *Log) Stop() error {
	return nil
}

func ModuleName() string {
	return "log"
}

func Methods() []string {
	return []string{
		"error",
		"warn",
		"info",
		"debug",
	}
}

func (l *Log) GetLoader(s *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"error": l.error(s.Name),
			"warn":  l.warn(s.Name),
			"info":  l.info(s.Name),
			"debug": l.debug(s.Name),
		}

		mod := luaState.SetFuncs(luaState.NewTable(), exports)

		luaState.Push(mod)
		return 1
	}
}

func (l *Log) error(scriptName string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		l.logger.Error(luaState.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}

func (l *Log) warn(scriptName string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		l.logger.Warn(luaState.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}

func (l *Log) info(scriptName string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		l.logger.Info(luaState.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}

func (l *Log) debug(scriptName string) lua.LGFunction {
	return func(luaState *lua.LState) int {
		l.logger.Debug(luaState.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}
