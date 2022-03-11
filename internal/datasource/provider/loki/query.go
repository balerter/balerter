package loki

import (
	"fmt"

	"github.com/balerter/balerter/internal/datasource/provider/loki/models"

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

func (m *Loki) doQuery(luaState *lua.LState) int {
	query, err := m.getQuery(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2
	}

	queryOptions := &queryOptions{}
	err = m.parseOptions(luaState, queryOptions)
	if err != nil {
		m.logger.Error("error parse query options", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error parse query options"))
		return 2
	}

	m.logger.Debug("call loki query", zap.String("name", m.name), zap.String("query", query), zap.Any("options", queryOptions))

	return m.do(luaState, m.sendQuery(query, queryOptions))
}

func (m *Loki) doRange(luaState *lua.LState) int {
	query, err := m.getQuery(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2 // nolint:gomnd // params count
	}

	rangeOptions := &rangeOptions{}
	err = m.parseOptions(luaState, rangeOptions)
	if err != nil {
		m.logger.Error("error parse range options", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error parse range options"))
		return 2 // nolint:gomnd // params count
	}

	m.logger.Debug("call loki query range", zap.String("name", m.name), zap.String("query", query), zap.Any("options", rangeOptions))

	return m.do(luaState, m.sendRange(query, rangeOptions))
}

func (m *Loki) do(luaState *lua.LState, u string) int {
	v, err := m.send(u)
	if err != nil {
		m.logger.Error("error send query to loki", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error send query to loki: " + err.Error()))
		return 2
	}

	if v == nil || v.Data.Result == nil {
		m.logger.Error("unexpected response from loki")
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("unexpected response from loki"))
		return 2
	}

	switch v.Data.Result.Type() {
	case models.ResultTypeStream:
		vv := v.Data.Result.(models.Streams)

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
		return 2
	}

	luaState.Push(lua.LNil)

	return 2
}
