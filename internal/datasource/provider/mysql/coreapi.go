package mysql

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func (m *MySQL) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	if method != "query" {
		return nil, http.StatusBadRequest, fmt.Errorf("unknown method: %q", method)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeout)
	defer ctxCancel()

	rows, err := m.db.QueryContext(ctx, string(body))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	cct, _ := rows.ColumnTypes()

	dest := make([]any, 0)

	for range cct {
		dest = append(dest, new([]byte))
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
			case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT":
				v, err := strconv.ParseInt(raw, 10, 64)
				if err != nil {
					m.logger.Error("error parse value", zap.String("value", raw), zap.Error(err))
					return nil, http.StatusInternalServerError, err
				}
				row[c.Name()] = v
			case "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE":
				v, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					m.logger.Error("error parse value", zap.String("value", raw), zap.Error(err))
					return nil, http.StatusInternalServerError, err
				}
				row[c.Name()] = v
			case "DATETIME":
				v, err := time.Parse("2006-01-02 15:04:05", raw)
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
