package alert

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// ModuleName returns the module name
func ModuleName() string {
	return "alert"
}

// Methods returns module methods
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

// Alert represents alert core module
type Alert struct {
	storage   corestorage.Alert
	chManager ChManager
	logger    *zap.Logger
}

// ChManager is an interface of channel manager
type ChManager interface {
	Send(a *alert.Alert, text string, options *alert.Options)
}

// New create new Alert core module
func New(storage corestorage.Alert, chManager ChManager, logger *zap.Logger) *Alert {
	a := &Alert{
		storage:   storage,
		chManager: chManager,
		logger:    logger,
	}

	return a
}

// Name returns the module name
func (a *Alert) Name() string {
	return ModuleName()
}

// GetLoader returns the lua loader
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

// Stop the module
func (a *Alert) Stop() error {
	return nil
}
