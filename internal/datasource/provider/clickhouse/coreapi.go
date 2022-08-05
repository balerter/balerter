package clickhouse

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/balerter/balerter/internal/datasource/converter"
	lua "github.com/yuin/gopher-lua"

	"go.uber.org/zap"
)

func (m *Clickhouse) CoreApiHandler(_ []string, body []byte) (any, int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeout)
	defer ctxCancel()

	rows, err := m.db.QueryContext(ctx, string(body))
	if err != nil {
		m.logger.Error("error clickhouse query", zap.ByteString("query", body), zap.Error(err))
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	cct, _ := rows.ColumnTypes()

	dest := make([]interface{}, 0)
	ffs := make([]func(v interface{}) lua.LValue, 0)

	for range cct {
		dest = append(dest, new([]byte))
		ffs = append(ffs, converter.FromDateBytes)
	}

	var result []map[string]any

	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			m.logger.Error("error scan", zap.Error(err))
			return nil, http.StatusInternalServerError, err
		}

		row := map[string]any{}

		for idx, c := range cct {
			raw := string(*dest[idx].(*[]byte))

			switch c.DatabaseTypeName() {
			case "Int64", "Int32", "Int16", "Int8", "Int":
				v, err := strconv.ParseInt(raw, 10, 64)
				if err != nil {
					m.logger.Error("error parse value", zap.String("value", raw), zap.Error(err))
					return nil, http.StatusInternalServerError, err
				}
				row[c.Name()] = v
			case "UInt64", "UInt32", "UInt16", "UInt8", "Uint":
				v, err := strconv.ParseUint(raw, 10, 64)
				if err != nil {
					m.logger.Error("error parse value", zap.String("value", raw), zap.Error(err))
					return nil, http.StatusInternalServerError, err
				}
				row[c.Name()] = v
			case "Float64", "Float32":
				v, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					m.logger.Error("error parse value", zap.String("value", raw), zap.Error(err))
					return nil, http.StatusInternalServerError, err
				}
				row[c.Name()] = v
			case "DateTime":
				v, err := time.Parse(time.RFC3339, raw)
				if err != nil {
					m.logger.Error("error parse value", zap.String("value", raw), zap.Error(err))
					return nil, http.StatusInternalServerError, err
				}
				row[c.Name()] = v

			default:
				row[c.Name()] = raw
			}
		}

		result = append(result, row)
	}

	return result, 0, nil
}
