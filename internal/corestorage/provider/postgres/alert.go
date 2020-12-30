package postgres

import (
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

func (p *PostgresAlert) GetOrNew(string) (*alert.Alert, error) {
	panic("not implemented")
	return nil, nil
}

func (p *PostgresAlert) All() ([]*alert.Alert, error) {
	panic("not implemented")
	return nil, nil
}

func (p *PostgresAlert) Store(a *alert.Alert) error {
	panic("not implemented")
	return nil
}

func (p *PostgresAlert) Get(string) (*alert.Alert, error) {
	panic("not implemented")
	return nil, nil
}
