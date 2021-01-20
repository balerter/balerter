package sql

import (
	"time"

	"github.com/balerter/balerter/internal/corestorage"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// SQL implements CoreStorage with the SQL as a storage backend
type SQL struct {
	name   string
	db     *sqlx.DB
	alerts *PostgresAlert
	kv     *PostgresKV
}

func New(name, driver, connectionString, tableAlerts, tableKV string, timeout time.Duration, logger *zap.Logger) (*SQL, error) {
	conn, err := sqlx.Connect(driver, connectionString)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	p := &SQL{
		name:   name,
		db:     conn,
		alerts: &PostgresAlert{db: conn, table: tableAlerts, timeout: timeout, logger: logger},
		kv:     &PostgresKV{db: conn, table: tableKV, timeout: timeout, logger: logger},
	}

	return p, nil
}

func (p *SQL) Name() string {
	return p.name
}

func (p *SQL) Stop() error {
	return p.db.Close()
}

func (p *SQL) KV() corestorage.KV {
	return p.kv
}

func (p *SQL) Alert() corestorage.Alert {
	return p.alerts
}
