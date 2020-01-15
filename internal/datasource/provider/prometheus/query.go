package prometheus

import (
	"fmt"
	"github.com/prometheus/common/model"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func (m *Prometheus) querySingle(L *lua.LState) int {
	query := strings.TrimSpace(L.Get(1).String())
	if query == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("query must be not empty"))
		return 2
	}

	m.logger.Debug("call prometheus query", zap.String("name", m.name), zap.String("query", query))

	v, err := m.do(query)
	if err != nil {
		m.logger.Error("error send query to prometheus", zap.Error(err))
		L.Push(lua.LNil)
		L.Push(lua.LString("error send query to prometheus: " + err.Error()))
		return 2
	}

	//fmt.Printf("\n\n\n%+v\n\n\n", v)

	switch v.Type() {
	case model.ValVector:
		vv := v.(model.Vector)
		if len(vv) == 0 {
			return m.luaError(L, err)
		}

		vvs := vv[0].Value.String()
		vl, err := strconv.ParseFloat(vvs, 64)
		if err != nil {
			return m.luaError(L, err)
		}

		L.Push(lua.LNumber(vl))
	default:
		L.Push(lua.LNil)
		return m.luaError(L, fmt.Errorf("unexpected prom model type"))
	}

	L.Push(lua.LNil)

	return 2
}

func (m *Prometheus) luaError(L *lua.LState, err error) int {
	m.logger.Debug("query error", zap.Error(err))
	L.Push(lua.LNil)
	L.Push(lua.LString(err.Error()))
	return 2
}
