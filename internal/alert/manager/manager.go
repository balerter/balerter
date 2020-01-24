package manager

import (
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
	Send(*message.Message) error
	SendSuccess(*message.Message) error
	SendError(*message.Message) error
}

type alertInfo struct {
	ScriptName string
	Count      int
}

type Manager struct {
	logger   *zap.Logger
	channels map[string]alertChannel
	activeMx sync.RWMutex
	active   map[string]*alertInfo
}

func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:   logger,
		channels: make(map[string]alertChannel),
		active:   make(map[string]*alertInfo),
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
				"on":  m.on(script),
				"off": m.off(script),
			}

			mod := L.SetFuncs(L.NewTable(), exports)

			L.Push(mod)
			return 1
		}
	}()
}
