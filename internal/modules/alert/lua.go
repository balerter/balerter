package alert

import (
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"

	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/modules"
)

// GetLoader returns the lua loader
func (a *Alert) GetLoader(j modules.Job) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"warn":    a.callFromLua(j.Script().Channels, j.Script().Escalate, alert.LevelWarn),
				"warning": a.callFromLua(j.Script().Channels, j.Script().Escalate, alert.LevelWarn),

				"error": a.callFromLua(j.Script().Channels, j.Script().Escalate, alert.LevelError),
				"fail":  a.callFromLua(j.Script().Channels, j.Script().Escalate, alert.LevelError),

				"success": a.callFromLua(j.Script().Channels, j.Script().Escalate, alert.LevelSuccess),
				"ok":      a.callFromLua(j.Script().Channels, j.Script().Escalate, alert.LevelSuccess),

				"get": a.get(),
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

func (a *Alert) getAlertData(luaState *lua.LState) (alertName, alertText string, options *alert.Options, err error) {
	options = alert.NewOptions()

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

	// fields
	fieldsOpts := alertOptions.RawGetString("fields")
	if fieldsOpts != lua.LNil {
		if fieldsOpts.Type() != lua.LTTable {
			err = fmt.Errorf("fields option must be a table")
			return
		}
		var fieldsErr error
		fieldsOpts.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			if value.Type() != lua.LTString {
				fieldsErr = fmt.Errorf("option key must be a string, %s", value.String())
				return
			}
			if value2.Type() != lua.LTString {
				fieldsErr = fmt.Errorf("option value must be a string, %s", value2.String())
				return
			}
			options.Fields[value.String()] = value2.String()
		})
		if fieldsErr != nil {
			err = fieldsErr
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

func (a *Alert) callFromLua(scriptChannels []string, escalate map[int][]string, alertLevel alert.Level) lua.LGFunction {
	return func(luaState *lua.LState) int {
		name, text, options, err := a.getAlertData(luaState)
		if err != nil {
			a.logger.Error("error get args", zap.Error(err))
			luaState.Push(lua.LString("error get arguments: " + err.Error()))
			return 1
		}

		_, _, errCall := a.call(name, text, scriptChannels, escalate, alertLevel, options)
		if errCall != nil {
			a.logger.Error("error update an alert", zap.Error(errCall))
			luaState.Push(lua.LString("error update an alert: " + errCall.Error()))
			return 1
		}

		return 0
	}
}

func (a *Alert) get() lua.LGFunction {
	return func(luaState *lua.LState) int {
		name := luaState.Get(1)
		if name.Type() != lua.LTString {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("alert name must be a string"))
			return 2
		}

		alrt, err := a.storage.Get(name.String())
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error get alert: " + err.Error()))
			return 2
		}

		t := alrt.MarshalLua()

		luaState.Push(t)

		return 1
	}
}
