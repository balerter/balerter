package alert

import (
	"fmt"
	"reflect"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"

	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/modules"
)

func (a *Alert) GetLoaderJS(j modules.Job) require.ModuleLoader {
	return func(runtime *goja.Runtime, object *goja.Object) {
		o := object.Get("exports").(*goja.Object)

		o.Set("warn", a.callFromJS(runtime, j.Script().Channels, j.Script().Escalate, alert.LevelWarn))
		o.Set("warning", a.callFromJS(runtime, j.Script().Channels, j.Script().Escalate, alert.LevelWarn))

		o.Set("error", a.callFromJS(runtime, j.Script().Channels, j.Script().Escalate, alert.LevelError))
		o.Set("fail", a.callFromJS(runtime, j.Script().Channels, j.Script().Escalate, alert.LevelError))

		o.Set("success", a.callFromJS(runtime, j.Script().Channels, j.Script().Escalate, alert.LevelSuccess))
		o.Set("ok", a.callFromJS(runtime, j.Script().Channels, j.Script().Escalate, alert.LevelSuccess))

		o.Set("get", a.getJS(runtime))
	}
}

func (a *Alert) getAlertDataJS(runtime *goja.Runtime, call goja.FunctionCall) (alertName, alertText string, options *alert.Options, err error) {
	options = alert.NewOptions()

	if len(call.Arguments) < 2 {
		return "", "", nil, fmt.Errorf("alert name and text must be defined")
	}

	name := call.Argument(0).String()
	text := call.Argument(1).String()

	if len(call.Arguments) == 2 {
		return name, text, options, nil
	}

	o := call.Argument(2)
	if o.ExportType() == nil {
		return "", "", nil, fmt.Errorf("options is null")
	}
	if o.ExportType().Kind() != reflect.Map {
		return "", "", nil, fmt.Errorf("options must be an object")
	}
	oo := o.ToObject(runtime)
	for _, key := range oo.Keys() {
		v := oo.Get(key)
		switch key {
		case "quiet":
			if v.ExportType() == nil || v.ExportType().Kind() != reflect.Bool {
				return "", "", nil, fmt.Errorf("option 'repeat' must be a boolean")
			}
			options.Quiet = v.ToBoolean()

		case "repeat", "resend":
			if v.ExportType() == nil || v.ExportType().Kind() != reflect.Int64 {
				return "", "", nil, fmt.Errorf("option 'repeat' must be an integer")
			}
			options.Repeat = int(v.ToInteger())
		default:
			return "", "", nil, fmt.Errorf("unknown options key: %s", key)
		}
	}

	return name, text, options, nil
}

func (a *Alert) callFromJS(runtime *goja.Runtime, scriptChannels []string, escalate map[int][]string, alertLevel alert.Level) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		name, text, options, err := a.getAlertDataJS(runtime, call)
		if err != nil {
			panic(runtime.ToValue("error parse alert data, " + err.Error()))
			return nil
		}

		_, _, errCall := a.call(name, text, scriptChannels, escalate, alertLevel, options)
		if errCall != nil {
			panic(runtime.ToValue("error update an alert, " + errCall.Error()))
			return nil
		}

		return nil
	}
}

func (a *Alert) getJS(runtime *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		name := call.Argument(0).String()

		alrt, err := a.storage.Get(name)
		if err != nil {
			panic(runtime.ToValue("error get alert '" + name + "', " + err.Error()))
			return nil
		}

		return runtime.ToValue(alrt)
	}
}
