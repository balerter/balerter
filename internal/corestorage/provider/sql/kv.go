package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrNoRow = errors.New("not found")
)

type PostgresKV struct {
	db      *sqlx.DB
	table   string
	timeout time.Duration
	logger  *zap.Logger
}

func (p *PostgresKV) All() (map[string]string, error) {
	query := fmt.Sprintf(`SELECT key, value FROM %s`, p.table)

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
	query := fmt.Sprintf(`INSERT INTO %s (key, value) VALUES ($1, $2) ON CONFLICT (key) DO NOTHING`, p.table)

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
	query := fmt.Sprintf(`SELECT value FROM %s WHERE key = $1`, p.table)

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
	query := fmt.Sprintf(`INSERT INTO %s (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = $2`, p.table)

	_, err := p.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("error sql query, %w", err)
	}

	return nil
}

func (p *PostgresKV) Delete(key string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE key = $1`, p.table)

	row := p.db.QueryRow(query, key)
	err := row.Err()
	if err != nil {
		return fmt.Errorf("error sql query, %w", err)
	}

	return nil
}
