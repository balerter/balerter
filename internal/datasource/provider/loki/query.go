package loki

import (
	"fmt"
	lokihttp "github.com/grafana/loki/pkg/loghttp"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (m *Loki) getQuery(luaState *lua.LState) (string, error) {
	queryV := luaState.Get(1)
	if queryV.Type() == lua.LTNil {
		return "", fmt.Errorf("query must be not empty")
	}
	query := queryV.String()
	if query == "" {
		return "", fmt.Errorf("query must be not empty")
	}
	return query, nil
}

func (m *Loki) doQuery(luaState *lua.LState) int { //nolint:dupl // not same code!
	query, err := m.getQuery(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2 //nolint:mnd
	}

	queryOptions, err := m.parseQueryOptions(luaState)
	if err != nil {
		m.logger.Error("error parse query options", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error parse query options"))
		return 2 //nolint:mnd
	}

	m.logger.Debug("call loki query", zap.String("name", m.name), zap.String("query", query), zap.Any("options", queryOptions))

	v, err := m.send(m.sendQuery(query, queryOptions))
	if err != nil {
		m.logger.Error("error send query to loki", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send query to loki: " + err.Error()))
		return 2 //nolint:mnd
	}

	return m.do(v, luaState)
}

func (m *Loki) doRange(luaState *lua.LState) int { //nolint:dupl // not same code!
	query, err := m.getQuery(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2 //nolint:mnd
	}

	rangeOptions, err := m.parseRangeOptions(luaState)
	if err != nil {
		m.logger.Error("error parse range options", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error parse range options"))
		return 2 //nolint:mnd
	}

	m.logger.Debug("call loki query range", zap.String("name", m.name), zap.String("query", query), zap.Any("options", rangeOptions))

	v, err := m.send(m.sendRange(query, rangeOptions))
	if err != nil {
		m.logger.Error("error send query to loki", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send query to loki: " + err.Error()))
		return 2 //nolint:mnd
	}

	return m.do(v, luaState)
}

func (m *Loki) do(v *lokihttp.QueryResponse, luaState *lua.LState) int {
	switch v.Data.Result.Type() {
	case lokihttp.ResultTypeStream:
		vv := v.Data.Result.(lokihttp.Streams)

		tbl := &lua.LTable{}
		for _, s := range vv {
			row := &lua.LTable{}
			labels := &lua.LTable{}
			for key, val := range s.Labels {
				labels.RawSet(lua.LString(key), lua.LString(val))
			}
			row.RawSet(lua.LString("labels"), labels)

			entries := &lua.LTable{}
			for _, e := range s.Entries {
				value := &lua.LTable{}
				value.RawSet(lua.LString("timestamp"), lua.LNumber(e.Timestamp.Unix()))
				value.RawSet(lua.LString("line"), lua.LString(e.Line))
				entries.Append(value)
			}
			row.RawSet(lua.LString("entries"), entries)
			tbl.Append(row)
		}

		luaState.Push(tbl)
	default:
		m.logger.Error("query error: unexpected loki model type")
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("query error: unexpected loki model type"))
		return 2 //nolint:mnd
	}

	luaState.Push(lua.LNil)

	return 2 //nolint:mnd
}
