package sql

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/metrics"
	"go.uber.org/zap"
	"time"
)

func (p *PostgresAlert) Update(name string, level alert.Level) (*alert.Alert, bool, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, false, fmt.Errorf("error start tx, %w", err)
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES ($1, $2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) ON CONFLICT (%s) DO NOTHING`,
		p.tableCfg.Table,
		p.tableCfg.Fields.Name,
		p.tableCfg.Fields.Level,
		p.tableCfg.Fields.Count,
		p.tableCfg.Fields.UpdatedAt,
		p.tableCfg.Fields.CreatedAt,
		p.tableCfg.Fields.Name,
	)

	res, err := tx.Exec(query, name, level)
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
		metrics.SetAlertLevel(name, level)
		return a, level != alert.LevelSuccess, nil
	}

	query = fmt.Sprintf(`SELECT %s, %s, %s, %s FROM %s WHERE %s = $1`,
		p.tableCfg.Fields.Level,
		p.tableCfg.Fields.Count,
		p.tableCfg.Fields.UpdatedAt,
		p.tableCfg.Fields.CreatedAt,
		p.tableCfg.Table,
		p.tableCfg.Fields.Name,
	)

	row := tx.QueryRow(query, name)

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
		query = fmt.Sprintf(`UPDATE %s SET %s = %s + 1, %s = CURRENT_TIMESTAMP WHERE %s = $1`,
			p.tableCfg.Table,
			p.tableCfg.Fields.Count,
			p.tableCfg.Fields.Count,
			p.tableCfg.Fields.UpdatedAt,
			p.tableCfg.Fields.Name,
		)

		_, err = tx.Exec(query, name)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				p.logger.Error("error rollback tx", zap.Error(err))
			}
			p.logger.Error("error update row", zap.Error(err))
			return nil, false, fmt.Errorf("error update row, %w", err)
		}
		a.Count++
		err = tx.Commit()
		if err == nil {
			metrics.SetAlertLevel(name, level)
		}
		return a, false, err
	}

	query = fmt.Sprintf(`UPDATE %s SET %s = $1, %s = 1, %s = CURRENT_TIMESTAMP WHERE %s = $2`,
		p.tableCfg.Table,
		p.tableCfg.Fields.Level,
		p.tableCfg.Fields.Count,
		p.tableCfg.Fields.UpdatedAt,
		p.tableCfg.Fields.Name,
	)

	_, err = tx.Exec(query, level, name)
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
	err = tx.Commit()
	if err == nil {
		metrics.SetAlertLevel(name, level)
	}
	return a, true, err
}
