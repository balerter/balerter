package postgres

import (
	"context"
	"fmt"

	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func (m *Postgres) query(luaState *lua.LState) int {
	q := luaState.Get(1).String()

	m.logger.Debug("call postgres query", zap.String("query", q))

	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeout)
	defer ctxCancel()

	rows, err := m.db.Query(ctx, q)
	if err != nil {
		m.logger.Error("error postgres query", zap.String("query", q), zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString(err.Error()))
		return 2
	}
	defer rows.Close()

	result := &lua.LTable{}

	for rows.Next() {
		values, errValues := rows.Values()
		if errValues != nil {
			m.logger.Error("error get values", zap.Error(errValues))
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error get values: " + errValues.Error()))
			return 2
		}

		row := &lua.LTable{}

		for idx, fd := range rows.FieldDescriptions() {
			row.RawSet(lua.LString(fd.Name), lua.LString(fmt.Sprintf("%v", values[idx])))
		}

		result.Append(row)
	}
	if err := rows.Err(); err != nil {
		m.logger.Error("error next", zap.Error(err))
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error next: " + err.Error()))
		return 2
	}

	luaState.Push(result)
	luaState.Push(lua.LNil)
	return 2
}
