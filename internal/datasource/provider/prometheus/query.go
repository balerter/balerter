package prometheus

import (
	"github.com/prometheus/common/model"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

type queryQueryOptions struct {
	Time string
}

func (m *Prometheus) doQuery(luaState *lua.LState) int {
	query := strings.TrimSpace(luaState.Get(1).String())
	if query == "" {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("query must be not empty"))
		return 2
	}

	options := luaState.Get(2)
	queryOptions := queryQueryOptions{}
	if options.Type() == lua.LTTable {
		err := gluamapper.Map(options.(*lua.LTable), &queryOptions)
		if err != nil {
			m.logger.Error("error decode query query options", zap.Error(err))
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error decode query query options"))
			return 2
		}
	}

	m.logger.Debug("call prometheus query", zap.String("name", m.name), zap.String("query", query))

	v, err := m.send(m.sendQuery(query, queryOptions))
	if err != nil {
		m.logger.Error("error send query to prometheus", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send query to prometheus: " + err.Error()))
		return 2
	}

	switch v.Type() {
	case model.ValVector:
		vv := v.(model.Vector)

		tbl := &lua.LTable{}
		for _, s := range vv {
			row := &lua.LTable{}
			metrics := &lua.LTable{}
			for key, val := range s.Metric {
				metrics.RawSet(lua.LString(key), lua.LString(val))
			}
			row.RawSet(lua.LString("metrics"), metrics)

			value := &lua.LTable{}
			value.RawSet(lua.LString("timestamp"), lua.LNumber(s.Timestamp.Unix()))
			value.RawSet(lua.LString("value"), lua.LNumber(s.Value))

			row.RawSet(lua.LString("value"), value)
			tbl.Append(row)
		}

		luaState.Push(tbl)
	default:
		m.logger.Debug("query error: unexpected prometheus model type")
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("query error: unexpected prometheus model type"))
		return 2
	}

	luaState.Push(lua.LNil)

	return 2
}

type queryRangeOptions struct {
	Start string
	End   string
	Step  string
}

func (m *Prometheus) doRange(luaState *lua.LState) int {
	query := strings.TrimSpace(luaState.Get(1).String())
	if query == "" {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("query must be not empty"))
		return 2
	}

	options := luaState.Get(2)
	rangeOptions := queryRangeOptions{}
	if options.Type() == lua.LTTable {
		err := gluamapper.Map(options.(*lua.LTable), &rangeOptions)
		if err != nil {
			m.logger.Error("error decode query range options", zap.Error(err))
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error decode query range options"))
			return 2
		}
	}

	m.logger.Debug("call prometheus query range", zap.String("name", m.name), zap.String("query", query))

	v, err := m.send(m.sendRange(query, rangeOptions))
	if err != nil {
		m.logger.Error("error send query to prometheus", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send query to prometheus: " + err.Error()))
		return 2
	}

	switch v.Type() {
	case model.ValMatrix:
		tbl := processValMatrixRange(v.(model.Matrix))
		luaState.Push(tbl)

	case model.ValVector:
		tbl := processValVectorRange(v.(model.Vector))
		luaState.Push(tbl)
	default:
		m.logger.Debug("query error: unexpected prometheus model type")
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("query error: unexpected prometheus model type"))
		return 2
	}

	luaState.Push(lua.LNil)

	return 2
}

func processValVectorRange(vv model.Vector) *lua.LTable {
	tbl := &lua.LTable{}
	for _, s := range vv {
		row := &lua.LTable{}
		metrics := &lua.LTable{}
		for key, val := range s.Metric {
			metrics.RawSet(lua.LString(key), lua.LString(val))
		}
		row.RawSet(lua.LString("metrics"), metrics)
		row.RawSet(lua.LString("value"), lua.LNumber(s.Value))
		tbl.Append(row)
	}

	return tbl
}

func processValMatrixRange(vv model.Matrix) *lua.LTable {
	tbl := &lua.LTable{}
	for _, s := range vv {
		row := &lua.LTable{}

		metrics := &lua.LTable{}
		for key, val := range s.Metric {
			metrics.RawSet(lua.LString(key), lua.LString(val))
		}

		values := &lua.LTable{}
		for _, val := range s.Values {
			value := &lua.LTable{}
			value.RawSet(lua.LString("timestamp"), lua.LNumber(val.Timestamp.Unix()))
			value.RawSet(lua.LString("value"), lua.LNumber(val.Value))

			values.Append(value)
		}

		row.RawSet(lua.LString("metrics"), metrics)
		row.RawSet(lua.LString("values"), values)
		tbl.Append(row)
	}

	return tbl
}
