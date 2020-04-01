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

func (l *Log) GetLoader(script *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"error": l.error(script.Name),
			"warn":  l.warn(script.Name),
			"info":  l.info(script.Name),
			"debug": l.debug(script.Name),
		}

		mod := L.SetFuncs(L.NewTable(), exports)

		L.Push(mod)
		return 1

	}
}

func (l *Log) error(scriptName string) lua.LGFunction {
	return func(L *lua.LState) int {
		l.logger.Error(L.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}

func (l *Log) warn(scriptName string) lua.LGFunction {
	return func(L *lua.LState) int {
		l.logger.Warn(L.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}

func (l *Log) info(scriptName string) lua.LGFunction {
	return func(L *lua.LState) int {
		l.logger.Info(L.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}

func (l *Log) debug(scriptName string) lua.LGFunction {
	return func(L *lua.LState) int {
		l.logger.Debug(L.Get(1).String(), zap.String("scriptName", scriptName))
		return 0
	}
}
