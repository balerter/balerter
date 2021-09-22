package alert

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

func (a *Alert) getAlertData(luaState *lua.LState) (alertName, alertText string, options *alert.Options, err error) {
	options = &alert.Options{}

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

	alertTextLua := luaState.Get(2) // nolint:gomnd // param position
	if alertTextLua.Type() == lua.LTNil {
		return
	}

	alertText = alertTextLua.String()

	alertOptionsLua := luaState.Get(3) // nolint:gomnd // param position
	if alertOptionsLua.Type() == lua.LTNil {
		return
	}

	if alertOptionsLua.Type() != lua.LTTable {
		err = fmt.Errorf("options must be a table")
		return
	}

	alertOptions := alertOptionsLua.(*lua.LTable)

	// parse options
	// channels
	channelsOpts := alertOptions.RawGetString("channels")
	if channelsOpts != lua.LNil {
		if channelsOpts.Type() != lua.LTTable {
			err = fmt.Errorf("channels options must be a table")
			return
		}
		var channelsErr error
		channelsOpts.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			if value2.Type() != lua.LTString {
				channelsErr = fmt.Errorf("wrong channel name %s", value2.String())
			}
			options.Channels = append(options.Channels, value2.String())
		})
		if channelsErr != nil {
			err = channelsErr
			return
		}
	}

	// quiet
	quietVal := alertOptions.RawGetString("quiet")
	if quietVal.Type() != lua.LTNil {
		if quietVal.Type() != lua.LTBool {
			err = fmt.Errorf("quiet must be a bool")
			return
		}
		options.Quiet = bool(quietVal.(lua.LBool))
	}

	// repeat & resend
	repeatVal := alertOptions.RawGetString("repeat")
	if repeatVal.Type() != lua.LTNil {
		if repeatVal.Type() != lua.LTNumber {
			err = fmt.Errorf("repeat must be a number")
			return
		}
		options.Repeat = int(repeatVal.(lua.LNumber))
	}

	resendVal := alertOptions.RawGetString("resend")
	if resendVal != lua.LNil {
		if alertOptionsLua.(*lua.LTable).RawGetString("repeat") != lua.LNil {
			err = fmt.Errorf("you must not use repeat and resend option together")
			return
		}

		if resendVal.Type() != lua.LTNumber {
			err = fmt.Errorf("resend must be a number")
			return
		}
		options.Repeat = int(resendVal.(lua.LNumber))
	}

	imageVal := alertOptions.RawGetString("image")
	if imageVal != lua.LNil {
		if imageVal.Type() != lua.LTString {
			err = fmt.Errorf("image must be a string")
			return
		}
		options.Image = imageVal.String()
	}

	return alertName, alertText, options, nil
}

func (a *Alert) call(scriptChannels []string, alertLevel alert.Level) lua.LGFunction {
	return func(luaState *lua.LState) int {
		name, text, options, err := a.getAlertData(luaState)
		if err != nil {
			a.logger.Error("error get args", zap.Error(err))
			luaState.Push(lua.LString("error get arguments: " + err.Error()))
			return 1
		}

		if len(options.Channels) == 0 {
			options.Channels = scriptChannels
		}

		updatedAlert, levelWasUpdated, err := a.storage.Update(name, alertLevel)
		if err != nil {
			a.logger.Error("error update an alert", zap.Error(err))
			luaState.Push(lua.LString("error update an alert: " + err.Error()))
			return 1
		}

		if levelWasUpdated {
			a.chManager.Send(updatedAlert, text, options)
		}

		return 0
	}
}
