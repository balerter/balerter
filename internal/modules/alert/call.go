package alert

import (
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func (a *Alert) getAlertData(luaState *lua.LState) (alertName, alertText string, options *alert2.Options, err error) {
	options = &alert2.Options{}

	alertNameLua := luaState.Get(1)
	if alertNameLua.Type() == lua.LTNil {
		err = fmt.Errorf("alert name must be provided")
		return
	}

	alertName = strings.TrimSpace(alertNameLua.String())
	if alertName == "" {
		err = fmt.Errorf("alert name must be not empty")
		return
	}

	alertTextLua := luaState.Get(2)
	if alertTextLua.Type() == lua.LTNil {
		return
	}

	alertText = alertTextLua.String()

	alertOptionsLua := luaState.Get(3)
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

func (a *Alert) call(s *script.Script, alertLevel alert2.Level) lua.LGFunction {
	return func(luaState *lua.LState) int {
		// TODO: wip

		//name, text, options, err := a.getAlertData(luaState)
		//if err != nil {
		//	a.logger.Error("error get args", zap.Error(err))
		//	luaState.Push(lua.LString("error get arguments: " + err.Error()))
		//	return 1
		//}
		//
		//if len(options.Channels) == 0 {
		//	options.Channels = s.Channels
		//}

		//err = a.storage.Update(name, alertLevel, text, options)
		//if err != nil {
		//	a.logger.Error("error call alert storage", zap.Error(err))
		//	luaState.Push(lua.LString("error call alert storage: " + err.Error()))
		//	return 1
		//}

		return 0
	}
}
