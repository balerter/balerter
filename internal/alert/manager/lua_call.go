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
	Fields []string
	Quiet  bool
}

func defaultOptions() options {
	return options{
		Fields: nil,
		Quiet:  false,
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

		_, _, _ = alertName, alertText, alertOptions

		m.alertsMx.Lock()
		defer m.alertsMx.Unlock()
		//info, ok := m.alerts[alertName]
		//if !ok {
		//	info = &alertInfo{
		//		Active:     false,
		//		ScriptName: s.Name,
		//	}
		//	m.alerts[alertName] = info
		//}
		//
		//if info.Active {
		//	m.sendSuccess(alertName, text, fields...)
		//}
		//
		//info.Active = false
		//
		//m.logger.Debug("call alert OFF", zap.String("alertName", alertName), zap.String("scriptName", s.Name))

		return 0
	}
}
