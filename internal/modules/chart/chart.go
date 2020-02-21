package chart

import (
	"bytes"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type DataItem struct {
	Timestamp float64
	Value     float64
}

type DataSeries struct {
	Color      string
	LineColor  string
	PointColor string
	Data       []DataItem
}

type Data struct {
	Series []DataSeries
}

type Chart struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Chart {
	l := &Chart{
		logger: logger,
	}

	return l
}

func (ch *Chart) Name() string {
	return "chart"
}

func (ch *Chart) Stop() error {
	return nil
}

func (ch *Chart) GetLoader(script *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"render": ch.render(script),
		}

		mod := L.SetFuncs(L.NewTable(), exports)

		L.Push(mod)
		return 1

	}
}

func (ch *Chart) render(_ *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		ch.logger.Info("Chart.Render")

		chartTitle := L.Get(1)
		if chartTitle.Type() == lua.LTNil {
			L.Push(lua.LNil)
			L.Push(lua.LString("title must be defined"))
			return 2
		}

		chartData := L.Get(2)
		if chartData.Type() != lua.LTTable {
			L.Push(lua.LNil)
			L.Push(lua.LString("chart data table must be defined"))
			return 2
		}

		data := &Data{}

		err := gluamapper.Map(chartData.(*lua.LTable), data)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("wrong chart data format, " + err.Error()))
			return 2
		}

		buf := bytes.NewBuffer([]byte{})

		err = ch.Render(chartTitle.String(), data, buf)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString("error render chart, " + err.Error()))
			return 2
		}

		L.Push(lua.LString(buf.String()))

		return 1
	}
}
