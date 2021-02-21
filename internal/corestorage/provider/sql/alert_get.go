package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"time"
)

func (p *PostgresAlert) Get(alertName string) (*alert.Alert, error) {
	row := p.db.QueryRow(fmt.Sprintf("SELECT id, level, count, last_change, start FROM %s WHERE id = $1", p.table), alertName)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error select alert, %w", err)
	}

	var l alert.Level
	var name string
	var level, count int
	var lastChange, start time.Time

	err := row.Scan(
		&name,
		&level,
		&count,
		&lastChange,
		&start,
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
	a.LastChange = lastChange
	a.Start = start

	return a, nil
}
