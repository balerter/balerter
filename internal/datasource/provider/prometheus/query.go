package prometheus

import (
	"fmt"
	"github.com/prometheus/common/model"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (m *Prometheus) getQuery(luaState *lua.LState) (string, error) {
	queryV := luaState.Get(1)
	if queryV.Type() != lua.LTString {
		return "", fmt.Errorf("query must be a string")
	}
	query := string(queryV.(lua.LString))
	if query == "" {
		return "", fmt.Errorf("query must be not empty")
	}
	return query, nil
}

func (m *Prometheus) doQuery(luaState *lua.LState) int {
	query, err := m.getQuery(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2
	}

	options, err := m.parseQueryOptions(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2
	}

	m.logger.Debug("call prometheus query", zap.String("name", m.name), zap.String("query", query))

	v, err := m.send(m.sendQuery(query, options))
	if err != nil {
		m.logger.Error("error send query to prometheus", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send query to prometheus: " + err.Error()))
		return 2
	}

	return m.resp(v, luaState)
}

func (m *Prometheus) doRange(luaState *lua.LState) int {
	query, err := m.getQuery(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2
	}

	options, err := m.parseRangeOptions(luaState)
	if err != nil {
		m.logger.Error("error decode query range options", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error decode query range options, " + err.Error()))
		return 2
	}

	m.logger.Debug("call prometheus query range", zap.String("name", m.name), zap.String("query", query))

	v, err := m.send(m.sendRange(query, options))
	if err != nil {
		m.logger.Error("error send query to prometheus", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send query to prometheus: " + err.Error()))
		return 2
	}

	return m.resp(v, luaState)
}

func (m *Prometheus) resp(v model.Value, luaState *lua.LState) int {
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
