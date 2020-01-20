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

type alertInfo struct {
	count int
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

func (m *Manager) getAlertName(L *lua.LState) (string, bool) {
	alertName := L.Get(1).String()
	alertName = strings.TrimSpace(alertName)

	if alertName == "" || alertName == "nil" {
		return "", false
	}

	return alertName, true
}
