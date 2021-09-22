package log

import (
	"github.com/balerter/balerter/internal/modules"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// Log represents the Log core module
type Log struct {
	logger *zap.Logger
}

// New creates new Log core module
func New(logger *zap.Logger) *Log {
	l := &Log{
		logger: logger,
	}

	return l
}

// Name returns the module name
func (l *Log) Name() string {
	return ModuleName()
}

// Stop the module
func (l *Log) Stop() error {
	return nil
}

// ModuleName returns the module name
func ModuleName() string {
	return "log"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"error",
		"warn",
		"info",
		"debug",
	}
}

// GetLoader returns the lua laoder
func (l *Log) GetLoader(j modules.Job) lua.LGFunction {
	return func(luaState *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"error": l.error(j.Script().Name),
			"warn":  l.warn(j.Script().Name),
			"info":  l.info(j.Script().Name),
			"debug": l.debug(j.Script().Name),
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
