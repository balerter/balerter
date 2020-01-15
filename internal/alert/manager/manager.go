package manager

import (
	"github.com/balerter/balerter/internal/alert/slack"
	"github.com/balerter/balerter/internal/config"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type alertChannel interface {
	Name() string
	Send(string, string) error
	SendSuccess(string, string) error
	SendError(string, string) error
}

type Manager struct {
	logger   *zap.Logger
	channels map[string]alertChannel
	activeMx sync.RWMutex
	active   map[string]int
}

func New(logger *zap.Logger) *Manager {
	m := &Manager{
		logger:   logger,
		channels: make(map[string]alertChannel),
		active:   make(map[string]int),
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

func (m *Manager) Loader() lua.LGFunction {
	return m.loader
}

func (m *Manager) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"on":  m.on,
		"off": m.off,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	//L.SetField(mod, "name", lua.LString("value"))

	L.Push(mod)
	return 1
}

func (m *Manager) on(L *lua.LState) int {
	alertName, ok := m.getAlertName(L)
	if !ok {
		return 0
	}

	alertText := L.Get(2).String()

	m.activeMx.Lock()
	defer m.activeMx.Unlock()
	if _, ok := m.active[alertName]; !ok {
		m.active[alertName] = 0
		m.sendError(alertName, alertText)
	}
	m.active[alertName]++

	m.logger.Debug("call alert ON", zap.String("alertName", alertName), zap.Int("count", m.active[alertName]))

	return 0
}

func (m *Manager) off(L *lua.LState) int {
	alertName, ok := m.getAlertName(L)
	if !ok {
		return 0
	}

	alertText := L.Get(2).String()

	m.activeMx.Lock()
	defer m.activeMx.Unlock()
	if _, ok := m.active[alertName]; ok {
		delete(m.active, alertName)
		m.sendSuccess(alertName, alertText)
	}

	m.logger.Debug("call alert OFF", zap.String("alertName", alertName))

	return 0
}

func (m *Manager) getAlertName(L *lua.LState) (string, bool) {
	alertName := L.Get(1).String()
	alertName = strings.TrimSpace(alertName)

	if alertName == "" {
		return "", false
	}

	return alertName, true
}

func (m *Manager) sendSuccess(alertName, message string) {
	for name, module := range m.channels {
		if err := module.SendSuccess(alertName, message); err != nil {
			m.logger.Error("error send message to channel", zap.String("name", name), zap.Error(err))
		}
	}
}

func (m *Manager) sendError(alertName, message string) {
	for name, module := range m.channels {
		if err := module.SendError(alertName, message); err != nil {
			m.logger.Error("error send message to channel", zap.String("name", name), zap.Error(err))
		}
	}
}
