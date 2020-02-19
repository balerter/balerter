package clickhouse

import (
	"context"
	"github.com/balerter/balerter/internal/datasource/converter"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

func (m *Clickhouse) query(L *lua.LState) int {

	q := L.Get(1).String()

	m.logger.Debug("call clickhouse query", zap.String("query", q))

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*3) // todo: timeout to settings
	defer ctxCancel()

	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		m.logger.Error("error clickhouse query", zap.String("query", q), zap.Error(err))
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
		case "UInt8", "UInt16", "UInt32", "UInt64":
			dest = append(dest, new(uint))
			ffs = append(ffs, converter.FromUInt)
		case "Int8", "Int16", "Int32", "Int64":
			dest = append(dest, new(int))
			ffs = append(ffs, converter.FromInt)
		case "Float32", "Float64":
			dest = append(dest, new(float64))
			ffs = append(ffs, converter.FromFloat64)
		case "String", "UUID":
			dest = append(dest, new(string))
			ffs = append(ffs, converter.FromString)
		case "Date":
			dest = append(dest, new(time.Time))
			ffs = append(ffs, converter.FromDate)
		case "DateTime":
			dest = append(dest, new(time.Time))
			ffs = append(ffs, converter.FromDateTime)
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
