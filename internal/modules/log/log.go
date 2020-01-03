package log

import (
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type Log struct {
	jobName string
	logger  *zap.Logger
}

func New(jobName string, logger *zap.Logger) lua.LGFunction {
	l := &Log{
		jobName: jobName,
		logger:  logger,
	}

	return l.Loader
}

func (l *Log) Loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"error": l.error,
		"warn":  l.warn,
		"info":  l.info,
		"debug": l.debug,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}

func (l *Log) error(L *lua.LState) int {
	l.logger.Error(L.Get(1).String())
	return 0
}
func (l *Log) warn(L *lua.LState) int {
	l.logger.Warn(L.Get(1).String())
	return 0
}
func (l *Log) info(L *lua.LState) int {
	l.logger.Info(L.Get(1).String())
	return 0
}
func (l *Log) debug(L *lua.LState) int {
	l.logger.Debug(L.Get(1).String())
	return 0
}
