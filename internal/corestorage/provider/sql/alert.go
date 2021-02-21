package sql

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type PostgresAlert struct {
	db      *sqlx.DB
	table   string
	timeout time.Duration
	logger  *zap.Logger
}

func (p *PostgresAlert) Update(name string, level alert.Level) (*alert.Alert, bool, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, false, fmt.Errorf("error start tx, %w", err)
	}

	res, err := tx.Exec(fmt.Sprintf(`INSERT INTO %s (id, level, count, last_change, start) VALUES ($1, $2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) ON CONFLICT (id) DO NOTHING`, p.table), name, level)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			p.logger.Error("error rollback tx", zap.Error(err))
		}
		p.logger.Error("error insert row", zap.Error(err))
		return nil, false, fmt.Errorf("error insert row, %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			p.logger.Error("error rollback tx", zap.Error(err))
		}
		return nil, false, fmt.Errorf("error get affected rows count, %w", err)
	}

	// if new alert
	if ra == 1 {
		err = tx.Commit()
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				p.logger.Error("error rollback tx", zap.Error(err))
			}
			return nil, false, fmt.Errorf("error commit tx, %w", err)
		}
		a := alert.New(name)
		a.Level = level
		return a, level != alert.LevelSuccess, nil
	}

	row := tx.QueryRow(fmt.Sprintf(`SELECT level, count, last_change, start FROM %s WHERE id = $1`, p.table), name)

	var l int
	var c int
	var lastChange time.Time
	var start time.Time

	err = row.Scan(&l, &c, &lastChange, &start)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			p.logger.Error("error rollback tx", zap.Error(err))
		}
		p.logger.Error("error scan row", zap.Error(err))
		return nil, false, fmt.Errorf("error scan row, %w", err)
	}

	currentLevel, err := alert.LevelFromInt(l)
	if err != nil {
		//err2 := tx.Commit()
		//if err2 != nil {
		//	err = fmt.Errorf("error commit tx, %w", err2)
		//}
		err2 := tx.Rollback()
		if err2 != nil {
			p.logger.Error("error rollback tx", zap.Error(err))
		}
		p.logger.Error("error convert level", zap.Error(err))
		return nil, false, fmt.Errorf("error convert level, %w", err)
	}

	a := alert.New(name)
	a.Level = currentLevel
	a.Count = c
	a.LastChange = lastChange
	a.Start = start

	// if level was not changed
	if currentLevel == level {
		_, err = tx.Exec(fmt.Sprintf(`UPDATE %s SET count = count + 1, last_change = CURRENT_TIMESTAMP WHERE id = $1`, p.table), name)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				p.logger.Error("error rollback tx", zap.Error(err))
			}
			p.logger.Error("error update row", zap.Error(err))
			return nil, false, fmt.Errorf("error update row, %w", err)
		}
		a.Count++
		return a, false, tx.Commit()
	}

	_, err = tx.Exec(fmt.Sprintf(`UPDATE %s SET level = $1, count = 1, last_change = CURRENT_TIMESTAMP WHERE id = $2`, p.table), level, name)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			p.logger.Error("error rollback tx", zap.Error(err))
		}
		p.logger.Error("error update row", zap.Error(err))
		return nil, true, fmt.Errorf("error update row, %w", err)
	}

	a.Count = 0
	a.Level = level

	return a, true, tx.Commit()
}
