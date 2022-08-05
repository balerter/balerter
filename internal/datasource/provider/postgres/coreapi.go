package postgres

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (m *Postgres) CoreApiHandler(method string, parts []string, params map[string]string, body []byte) (any, int, error) {
	if method != "query" {
		return nil, http.StatusBadRequest, fmt.Errorf("unknown method %q", method)
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	rows, err := m.db.Query(ctx, string(body))
	if err != nil {
		m.logger.Error("error postgres query", zap.ByteString("query", body), zap.Error(err))
		return nil, http.StatusBadRequest, err
	}
	defer rows.Close()

	var result []map[string]any

	for rows.Next() {
		values, errValues := rows.Values()
		if errValues != nil {
			m.logger.Error("error get rows values", zap.Error(errValues))
			return nil, http.StatusInternalServerError, errValues
		}

		row := map[string]any{}

		for idx, fd := range rows.FieldDescriptions() {
			row[string(fd.Name)] = values[idx]
		}

		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		m.logger.Error("rows error", zap.Error(err))
		return nil, http.StatusInternalServerError, err
	}

	return result, 0, nil
}
