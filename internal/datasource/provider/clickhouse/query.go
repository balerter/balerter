package clickhouse

import (
	"context"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func (m *Clickhouse) query(L *lua.LState) int {

	q := L.Get(1).String()

	m.logger.Debug("call clickhouse query", zap.String("query", q))

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*3)
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
	//cc, _ := rows.Columns()

	dest := make([]interface{}, 0)
	ffs := make([]func(v interface{}) string, 0)

	for _, c := range cct {
		switch c.DatabaseTypeName() {
		case "Float64":
			dest = append(dest, new(float64))
			ffs = append(ffs, convertToString_Float64)
		case "Date":
			dest = append(dest, new(time.Time))
			ffs = append(ffs, convertToString_Date)
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
			row.RawSet(lua.LString(c.Name()), lua.LString(v))
		}

		result.Append(row)
	}

	L.Push(result)
	L.Push(lua.LNil)
	return 2
}

func convertToString_Float64(v interface{}) string {
	return strconv.FormatFloat(*(v.(*float64)), 'f', -1, 64)
}

func convertToString_Date(v interface{}) string {
	return v.(*time.Time).Format("2006-01-02")
}
