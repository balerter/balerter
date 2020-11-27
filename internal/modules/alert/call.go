package alert

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

func (a *Alert) getAlertData(luaState *lua.LState) (alertName, alertText string, options *alert.Options, err error) {
	options = &alert.Options{}

	alertNameLua := luaState.Get(1) //nolint:mnd
	if alertNameLua.Type() == lua.LTNil {
		err = fmt.Errorf("alert name must be provided")
		return
	}

	alertName = strings.TrimSpace(alertNameLua.String())
	if alertName == "" {
		err = fmt.Errorf("alert name must be not empty")
		return
	}

	alertTextLua := luaState.Get(2) //nolint:mnd
	if alertTextLua.Type() == lua.LTNil {
		return
	}

	alertText = alertTextLua.String()

	alertOptionsLua := luaState.Get(3) //nolint:mnd
	if alertOptionsLua.Type() == lua.LTNil {
		return
	}

	if alertOptionsLua.Type() != lua.LTTable {
		err = fmt.Errorf("options must be a table")
		return
	}

	err = gluamapper.Map(alertOptionsLua.(*lua.LTable), &options)
	if err != nil {
		err = fmt.Errorf("wrong options format: %v", err)
		return
	}

	return alertName, alertText, options, nil
}

func (a *Alert) call(s *script.Script, alertLevel alert.Level) lua.LGFunction {
	return func(luaState *lua.LState) int {
		name, text, options, err := a.getAlertData(luaState)
		if err != nil {
			a.logger.Error("error get args", zap.Error(err))
			luaState.Push(lua.LString("error get arguments: " + err.Error()))
			return 1
		}

		if len(options.Channels) == 0 {
			options.Channels = s.Channels
		}

		err = a.manager.Call(name, alertLevel, text, options)
		if err != nil {
			a.logger.Error("error call alert manager", zap.Error(err))
			luaState.Push(lua.LString("error call alert manager: " + err.Error()))
			return 1
		}

		return 0
	}
}
