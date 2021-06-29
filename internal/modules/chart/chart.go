package chart

import (
	"bytes"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// ModuleName returns the module name
func ModuleName() string {
	return "chart"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"render",
	}
}

// DataItem represents data item
type DataItem struct {
	Timestamp float64
	Value     float64
}

// DataSeries represents data series
type DataSeries struct {
	Color      string
	LineColor  string
	PointColor string
	Data       []DataItem
}

// Data represents chart data
type Data struct {
	Title  string
	Series []DataSeries
}

// Chart represent Chart core module
type Chart struct {
	logger *zap.Logger
}

// New creates new Chart core module
func New(logger *zap.Logger) *Chart {
	l := &Chart{
		logger: logger,
	}

	return l
}

// Name returns the module name
func (ch *Chart) Name() string {
	return ModuleName()
}

// Stop the module
func (ch *Chart) Stop() error {
	return nil
}

// GetLoader returns the lua loader
func (ch *Chart) GetLoader(s *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"render": ch.render(s),
		}

		mod := luaState.SetFuncs(luaState.NewTable(), exports)

		luaState.Push(mod)
		return 1
	}
}

func (ch *Chart) render(_ *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		ch.logger.Debug("Chart.Render")

		chartTitle := luaState.Get(1)
		if chartTitle.Type() == lua.LTNil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("title must be defined"))
			return 2
		}

		chartData := luaState.Get(2)
		if chartData.Type() != lua.LTTable {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("chart data table must be defined"))
			return 2
		}

		data := &Data{}

		err := gluamapper.Map(chartData.(*lua.LTable), data)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("wrong chart data format, " + err.Error()))
			return 2
		}

		buf := bytes.NewBuffer([]byte{})

		err = ch.Render(chartTitle.String(), data, buf)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error render chart, " + err.Error()))
			return 2
		}

		luaState.Push(lua.LString(buf.String()))

		return 1
	}
}
