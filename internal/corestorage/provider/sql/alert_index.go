package sql

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Index is an implementation of the storage interface
func (p *PostgresAlert) Index(levels []alert.Level) (alert.Alerts, error) {
	query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s FROM %s",
		p.tableCfg.Fields.Name,
		p.tableCfg.Fields.Level,
		p.tableCfg.Fields.Count,
		p.tableCfg.Fields.UpdatedAt,
		p.tableCfg.Fields.CreatedAt,
		p.tableCfg.Table,
	)

	if len(levels) > 0 {
		var ll []string
		for _, l := range levels {
			ll = append(ll, l.NumString())
		}
		query += fmt.Sprintf(" WHERE %s IN (%s)", p.tableCfg.Fields.Level, strings.Join(ll, ","))
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
