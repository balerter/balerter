package alert

import (
	alert2 "github.com/balerter/balerter/internal/alert"
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
	storage corestorage.Alert
	logger  *zap.Logger
}

func New(storage corestorage.Alert, logger *zap.Logger) *Alert {
	a := &Alert{
		storage: storage,
		logger:  logger,
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
				"warn":    a.call(s, alert2.LevelWarn),
				"warning": a.call(s, alert2.LevelWarn),

				"error": a.call(s, alert2.LevelError),
				"fail":  a.call(s, alert2.LevelError),

				"success": a.call(s, alert2.LevelSuccess),
				"ok":      a.call(s, alert2.LevelSuccess),

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
