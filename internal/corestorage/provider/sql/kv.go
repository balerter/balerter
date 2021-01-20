package sql

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type PostgresKV struct {
	db      *sqlx.DB
	table   string
	timeout time.Duration
	logger  *zap.Logger
}

func (p *PostgresKV) All() (map[string]string, error) {
	panic("not implemented")
	return nil, nil
}

func (p *PostgresKV) Put(string, string) error {
	panic("not implemented")
	return nil
}

func (p *PostgresKV) Get(string) (string, error) {
	panic("not implemented")
	return "", nil
}

func (p *PostgresKV) Upsert(string, string) error {
	panic("not implemented")
	return nil
}

func (p *PostgresKV) Delete(string) error {
	panic("not implemented")
	return nil
}
