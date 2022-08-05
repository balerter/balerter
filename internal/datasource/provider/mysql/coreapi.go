package mysql

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/datasource/converter"
	lua "github.com/yuin/gopher-lua"
	"net/http"
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

	var result []map[string]any

	cct, _ := rows.ColumnTypes()

	dest := make([]interface{}, 0)
	ffs := make([]func(v interface{}) lua.LValue, 0)

	for range cct {
		dest = append(dest, new([]byte))
		ffs = append(ffs, converter.FromDateBytes)
	}

	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		row := map[string]any{}

		for idx, c := range cct {
			v := ffs[idx](dest[idx])
			row[c.Name()] = v
		}

		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return result, 0, nil
}
