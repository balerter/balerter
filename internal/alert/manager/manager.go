package manager

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/alert/slack"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"sync"
)

type alertChannel interface {
	Name() string
	Send(alert.Level, *message.Message) error
}

type Manager struct {
	logger   *zap.Logger
	channels map[string]alertChannel

	alertsMx sync.RWMutex
	alerts   map[string]*alert.Alert
}

func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:   logger,
		channels: make(map[string]alertChannel),
		alerts:   make(map[string]*alert.Alert),
	}

	return m
}

func (m *Manager) Init(cfg config.Channels) error {
	for _, configWebHook := range cfg.Slack {
		module, err := slack.New(configWebHook, m.logger)
		if err != nil {
			return err
		}

		m.channels[module.Name()] = module
	}

	return nil
}

func (m *Manager) Name() string {
	return "alert"
}

func (m *Manager) Stop() error {
	return nil
}

func (m *Manager) GetLoader(script *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(L *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"warn": m.luaCall(script, alert.LevelWarn),
				"info": m.luaCall(script, alert.LevelInfo),

				"error": m.luaCall(script, alert.LevelError),
				"on":    m.luaCall(script, alert.LevelError),
				"fail":  m.luaCall(script, alert.LevelError),

				"success": m.luaCall(script, alert.LevelSuccess),
				"off":     m.luaCall(script, alert.LevelSuccess),
				"ok":      m.luaCall(script, alert.LevelSuccess),
			}

			mod := L.SetFuncs(L.NewTable(), exports)

			L.Push(mod)
			return 1
		}
	}()
}
