package sql

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (p *PostgresAlert) Index(levels []alert.Level) (alert.Alerts, error) {
	query := fmt.Sprintf("SELECT id, level, count, last_change, start FROM %s", p.table)
	if len(levels) > 0 {
		var ll []string
		for _, l := range levels {
			ll = append(ll, l.NumString())
		}
		query += fmt.Sprintf(" WHERE level IN (%s)", strings.Join(ll, ","))
	}

	p.logger.Debug("select alerts index", zap.String("query", query))

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error select rows, %w", err)
	}

	result := make([]*alert.Alert, 0)

	var l alert.Level
	var name string
	var level, count int
	var lastChange, start time.Time

	for rows.Next() {

		err = rows.Scan(
			&name,
			&level,
			&count,
			&lastChange,
			&start,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan result, %w", err)
		}

		l, err = alert.LevelFromInt(level)
		if err != nil {
			return nil, fmt.Errorf("error parse level %d for alert %s, %w", level, name, err)
		}

		a := alert.New(name)
		a.Count = count
		a.Level = l
		a.LastChange = lastChange
		a.Start = start

		result = append(result, a)
	}

	return result, nil
}
