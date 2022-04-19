package alert

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/modules"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

//go:generate moq -out ch_manager_mock.go -skip-ensure -fmt goimports . chManager

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
	chManager chManager
	logger    *zap.Logger
}

// chManager is an interface of channel manager
type chManager interface {
	Send(a *alert.Alert, text string, options *alert.Options)
}

// New create new Alert core module
func New(storage corestorage.Alert, chManager chManager, logger *zap.Logger) *Alert {
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
func (a *Alert) GetLoader(j modules.Job) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"warn":    a.call(j.Script().Channels, j.Script().Escalate, alert.LevelWarn),
				"warning": a.call(j.Script().Channels, j.Script().Escalate, alert.LevelWarn),

				"error": a.call(j.Script().Channels, j.Script().Escalate, alert.LevelError),
				"fail":  a.call(j.Script().Channels, j.Script().Escalate, alert.LevelError),

				"success": a.call(j.Script().Channels, j.Script().Escalate, alert.LevelSuccess),
				"ok":      a.call(j.Script().Channels, j.Script().Escalate, alert.LevelSuccess),

				"get": a.get(),
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
