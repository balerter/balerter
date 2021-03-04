package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrNoRow = errors.New("not found")
)

type PostgresKV struct {
	db       *sqlx.DB
	tableCfg tables.TableKV
	timeout  time.Duration
	logger   *zap.Logger
}

func (p *PostgresKV) All() (map[string]string, error) {
	query := fmt.Sprintf(`SELECT %s, %s FROM %s`, p.tableCfg.Fields.Key, p.tableCfg.Fields.Value, p.tableCfg.Table)

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error sql query, %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error rows, %w", err)
	}

	result := make(map[string]string)

	var key, value string

	for rows.Next() {
		err = rows.Scan(
			&key,
			&value,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan result, %w", err)
		}
		result[key] = value
	}

	return result, nil
}

func (p *PostgresKV) Put(key, value string) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES ($1, $2) ON CONFLICT (%s) DO NOTHING`, p.tableCfg.Table, p.tableCfg.Fields.Key, p.tableCfg.Fields.Value, p.tableCfg.Fields.Key)

	res, err := p.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("error sql query, %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error get affected rows, %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("key already exists")
	}

	return nil
}

func (p *PostgresKV) Get(key string) (string, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE key = $1`, p.tableCfg.Fields.Value, p.tableCfg.Table)

	row := p.db.QueryRow(query, key)
	err := row.Err()
	if err != nil {
		return "", fmt.Errorf("error sql query, %w", err)
	}

	var value string

	err = row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRow
		}
		return "", fmt.Errorf("error scan result, %w", err)
	}

	return value, nil
}

func (p *PostgresKV) Upsert(key, value string) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES ($1, $2) ON CONFLICT (%s) DO UPDATE SET value = $2`, p.tableCfg.Table, p.tableCfg.Fields.Key, p.tableCfg.Fields.Value, p.tableCfg.Fields.Key)

	_, err := p.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("error sql query, %w", err)
	}

	return nil
}

func (p *PostgresKV) Delete(key string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, p.tableCfg.Table, p.tableCfg.Fields.Key)

	row := p.db.QueryRow(query, key)
	err := row.Err()
	if err != nil {
		return fmt.Errorf("error sql query, %w", err)
	}

	return nil
}
