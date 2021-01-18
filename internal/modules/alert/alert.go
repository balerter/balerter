package alert

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func ModuleName() string {
	return "alert"
}

func Methods() []string {
	return []string{
		"warn",
		"warning",
		"error",
		"fail",
		"success",
		"ok",
		"get",
	}
}

type Alert struct {
	storage   corestorage.Alert
	chManager ChManager
	logger    *zap.Logger
}

type ChManager interface {
	Send(a *alert.Alert, text string, options *alert.Options)
}

func New(storage corestorage.Alert, chManager ChManager, logger *zap.Logger) *Alert {
	a := &Alert{
		storage:   storage,
		chManager: chManager,
		logger:    logger,
	}

	return a
}

func (a *Alert) Name() string {
	return ModuleName()
}

func (a *Alert) GetLoader(s *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"warn":    a.call(s, alert.LevelWarn),
				"warning": a.call(s, alert.LevelWarn),

				"error": a.call(s, alert.LevelError),
				"fail":  a.call(s, alert.LevelError),

				"success": a.call(s, alert.LevelSuccess),
				"ok":      a.call(s, alert.LevelSuccess),

				"get": a.get(s),
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

func (a *Alert) Stop() error {
	return nil
}
