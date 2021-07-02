package sql

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// PostgresAlert represent Postgres implementation for Alert storage
type PostgresAlert struct {
	db       *sqlx.DB
	tableCfg tables.TableAlerts
	timeout  time.Duration
	logger   *zap.Logger
}

func (p *PostgresAlert) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS %s
(
	%s varchar not null,
	%s integer default 0 not null,
	%s integer default 0,
	%s timestamp default CURRENT_TIMESTAMP,
	%s timestamp default CURRENT_TIMESTAMP
);
`

	query = fmt.Sprintf(query,
		p.tableCfg.Table,
		p.tableCfg.Fields.Name,
		p.tableCfg.Fields.Level,
		p.tableCfg.Fields.Count,
		p.tableCfg.Fields.UpdatedAt,
		p.tableCfg.Fields.CreatedAt,
	)

	_, err := p.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
