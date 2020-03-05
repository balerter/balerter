package mysql

import (
	"context"
	"database/sql"
	"github.com/balerter/balerter/internal/datasource/converter"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

func (m *MySQL) query(L *lua.LState) int {

	q := L.Get(1).String()

	m.logger.Debug("call mysql query", zap.String("query", q))

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*3) // todo: timeout to settings
	defer ctxCancel()

	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		m.logger.Error("error mysql query", zap.String("query", q), zap.Error(err))
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer rows.Close()

	cct, _ := rows.ColumnTypes()

	dest := make([]interface{}, 0)
	ffs := make([]func(v interface{}) lua.LValue, 0)

	for _, c := range cct {
		switch c.DatabaseTypeName() {
		case "FLOAT":
			dest = append(dest, new(sql.NullFloat64))
			ffs = append(ffs, converter.FromFloat64)
		case "BLOB", "CHAR", "JSON", "TEXT", "VARCHAR":
			dest = append(dest, new(sql.NullString))
			ffs = append(ffs, converter.FromString)
		case "BIGINT", "BINARY", "TINYINT", "DECIMAL", "DOUBLE", "INT", "MEDIUMINT", "SMALLINT":
			dest = append(dest, new(sql.NullInt64))
			ffs = append(ffs, converter.FromInt)
		case "DATE", "DATETIME", "TIME", "TIMESTAMP", "YEAR":
			dest = append(dest, new([]byte))
			ffs = append(ffs, converter.FromDateBytes)
		case "BIT", "ENUM", "GEOMETRY":
			m.logger.Error("error scan type", zap.String("typename", c.DatabaseTypeName()))
			L.Push(lua.LNil)
			L.Push(lua.LString("unsupported database type"))
			return 2
		default:
			m.logger.Error("error scan type", zap.String("typename", c.DatabaseTypeName()))
			L.Push(lua.LNil)
			L.Push(lua.LString("error database type"))
			return 2
		}
	}

	result := &lua.LTable{}

	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			m.logger.Error("error scan", zap.Error(err))
			L.Push(lua.LNil)
			L.Push(lua.LString("error scan: " + err.Error()))
			return 2
		}

		row := &lua.LTable{}

		for idx, c := range cct {
			v := ffs[idx](dest[idx])
			row.RawSet(lua.LString(c.Name()), v)
		}

		result.Append(row)
	}

	L.Push(result)
	L.Push(lua.LNil)
	return 2
}
