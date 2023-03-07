package log

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"

	"github.com/balerter/balerter/internal/modules"
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

func (l *Log) GetLoaderJS(j modules.Job) require.ModuleLoader {
	return func(runtime *goja.Runtime, object *goja.Object) {
		o := object.Get("exports").ToObject(runtime)
		o.Set("error", l.jsLog(j.Script().Name, l.logger.Error))
		o.Set("warn", l.jsLog(j.Script().Name, l.logger.Warn))
		o.Set("info", l.jsLog(j.Script().Name, l.logger.Info))
		o.Set("debug", l.jsLog(j.Script().Name, l.logger.Debug))
	}
}

func (l *Log) jsLog(name string, fn func(msg string, fields ...zap.Field)) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			return nil
		}

		var fields []string

		for i := 1; i < len(call.Arguments); i++ {
			fields = append(fields, call.Argument(i).String())
		}

		fn(call.Argument(0).String(), zap.String("scriptName", name), zap.Strings("fields", fields))
		return nil
	}
}

// GetLoader returns the lua loader
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
