package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

type options struct {
	Fields   []string
	Channels []string
	Quiet    bool
	Repeat   int
}

func defaultOptions() options {
	return options{
		Fields: nil,
		Quiet:  false,
		Repeat: 0,
	}
}

func (m *Manager) getAlertData(L *lua.LState) (alertName string, alertText string, alertOptions options, err error) {
	alertOptions = defaultOptions()

	alertNameLua := L.Get(1)
	if alertNameLua.Type() == lua.LTNil {
		err = fmt.Errorf("alert name must be provided")
		return
	}

	alertName = strings.TrimSpace(alertNameLua.String())
	if alertName == "" {
		err = fmt.Errorf("alert name must be not empty")
		return
	}

	alertTextLua := L.Get(2)
	if alertTextLua.Type() == lua.LTNil {
		return
	}

	alertText = alertTextLua.String()

	alertOptionsLua := L.Get(3)
	if alertOptionsLua.Type() == lua.LTNil {
		return
	}

	if alertOptionsLua.Type() != lua.LTTable {
		err = fmt.Errorf("options must be a table")
		return
	}

	err = gluamapper.Map(alertOptionsLua.(*lua.LTable), &alertOptions)
	if err != nil {
		err = fmt.Errorf("wrong options format: %v", err)
		return
	}

	return
}

func (m *Manager) luaCall(s *script.Script, alertLevel alert.Level) lua.LGFunction {
	return func(L *lua.LState) int {
		alertName, alertText, alertOptions, err := m.getAlertData(L)
		if err != nil {
			m.logger.Error("error get args", zap.Error(err))
			L.Push(lua.LString("error get arguments: " + err.Error()))
			return 1
		}

		m.logger.Debug("call alert luaCall", zap.String("alertName", alertName), zap.String("scriptName", s.Name), zap.String("alertText", alertText), zap.Int("alertLevel", int(alertLevel)), zap.Any("alertOptions", alertOptions))

		m.alertsMx.Lock()
		a, ok := m.alerts[alertName]
		if !ok {
			a = alert.New()
			m.alerts[alertName] = a
		}
		m.alertsMx.Unlock()

		// if a level equals Info, always send the message (but only if not a quiet options)
		if alertLevel == alert.LevelInfo {
			if !alertOptions.Quiet {
				m.Send(alertLevel, alertName, alertText, alertOptions.Channels, alertOptions.Fields)
			}
			return 0
		}

		if a.Level() == alertLevel {
			a.Inc()

			if alertOptions.Repeat > 0 && a.Count()%alertOptions.Repeat == 0 {
				m.Send(alertLevel, alertName, alertText, alertOptions.Channels, alertOptions.Fields)
			}

			return 0
		}

		a.UpdateLevel(alertLevel)

		m.Send(alertLevel, alertName, alertText, alertOptions.Channels, alertOptions.Fields)

		return 0
	}
}
