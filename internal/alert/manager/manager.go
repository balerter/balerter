package manager

import (
	"fmt"

	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/alert/provider/email"
	"github.com/balerter/balerter/internal/alert/provider/notify"
	"github.com/balerter/balerter/internal/alert/provider/slack"
	"github.com/balerter/balerter/internal/alert/provider/syslog"
	"github.com/balerter/balerter/internal/alert/provider/telegram"
	"github.com/balerter/balerter/internal/config"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type alertChannel interface {
	Name() string
	Send(*message.Message) error
}

// Manager represents the Alert manager struct
type Manager struct {
	logger   *zap.Logger
	channels map[string]alertChannel

	engine coreStorage.CoreStorage
}

// New returns new Alert manager instance
func New(engine coreStorage.CoreStorage, logger *zap.Logger) *Manager {
	m := &Manager{
		logger:   logger,
		engine:   engine,
		channels: make(map[string]alertChannel),
	}

	return m
}

// Init the Alert manager
func (m *Manager) Init(cfg config.Channels) error {
	for _, cfg := range cfg.Email {
		module, err := email.New(cfg, m.logger)
		if err != nil {
			return fmt.Errorf("error init email channel %s, %w", cfg.Name, err)
		}

		m.channels[module.Name()] = module
	}

	for _, cfg := range cfg.Slack {
		module, err := slack.New(cfg, m.logger)
		if err != nil {
			return fmt.Errorf("error init slack channel %s, %w", cfg.Name, err)
		}

		m.channels[module.Name()] = module
	}

	for _, cfg := range cfg.Telegram {
		module, err := telegram.New(cfg, m.logger)
		if err != nil {
			return fmt.Errorf("error init telegram channel %s, %w", cfg.Name, err)
		}

		m.channels[module.Name()] = module
	}

	for _, cfg := range cfg.Syslog {
		module, err := syslog.New(&cfg, m.logger)
		if err != nil {
			return fmt.Errorf("error init syslog channel %s, %w", cfg.Name, err)
		}

		m.channels[module.Name()] = module
	}

	for _, cfg := range cfg.Notify {
		module, err := notify.New(cfg, m.logger)
		if err != nil {
			return fmt.Errorf("error init syslog channel %s, %w", cfg.Name, err)
		}

		m.channels[module.Name()] = module
	}

	return nil
}

// Name returns the alert manager name
func (m *Manager) Name() string {
	return ModuleName()
}

// Stop the alert manager
func (m *Manager) Stop() error {
	return nil
}

// ModuleName returns a lua module name for the alert manager
func ModuleName() string {
	return "alert"
}

// Methods returns the list of methods of the alert manager
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

// GetLoader implements Modules/Module.GetLoader
func (m *Manager) GetLoader(s *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"warn":    m.luaCall(s, alert.LevelWarn),
				"warning": m.luaCall(s, alert.LevelWarn),

				"error": m.luaCall(s, alert.LevelError),
				"fail":  m.luaCall(s, alert.LevelError),

				"success": m.luaCall(s, alert.LevelSuccess),
				"ok":      m.luaCall(s, alert.LevelSuccess),

				"get": m.get(s),
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}
