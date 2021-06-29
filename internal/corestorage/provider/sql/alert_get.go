package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"time"
)

// Get is an implementation of the storage interface
func (p *PostgresAlert) Get(alertName string) (*alert.Alert, error) {
	query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s FROM %s WHERE %s = $1",
		p.tableCfg.Fields.Name,
		p.tableCfg.Fields.Level,
		p.tableCfg.Fields.Count,
		p.tableCfg.Fields.UpdatedAt,
		p.tableCfg.Fields.CreatedAt,
		p.tableCfg.Table,
		p.tableCfg.Fields.Name,
	)

	row := p.db.QueryRow(query, alertName)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error select alert, %w", err)
	}

	var l alert.Level
	var name string
	var level, count int
	var updatedAt, createdAt time.Time

	err := row.Scan(
		&name,
		&level,
		&count,
		&updatedAt,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error scan result, %w", err)
	}

	l, err = alert.LevelFromInt(level)
	if err != nil {
		return nil, fmt.Errorf("error parse level %d for alert %s, %w", level, name, err)
	}

	a := alert.New(name)
	a.Level = l
	a.Count = count
	a.LastChange = updatedAt
	a.Start = createdAt

	return a, nil
}
