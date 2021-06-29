package sql

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"time"

	"github.com/balerter/balerter/internal/corestorage"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // import sqlite driver
	"go.uber.org/zap"
)

// SQL implements CoreStorage with the SQL as a storage backend
type SQL struct {
	name   string
	db     *sqlx.DB
	alerts *PostgresAlert
	kv     *PostgresKV
}

// New creates new SQL storage provider
func New(name, driver, connectionString string, alertsCfg tables.TableAlerts, kvCfg tables.TableKV, timeout time.Duration, logger *zap.Logger) (*SQL, error) {
	conn, err := sqlx.Connect(driver, connectionString)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		err2 := conn.Close()
		if err2 != nil {
			return nil, fmt.Errorf("error close sql connection after wrong ping %v, %w", err2, err)
		}
		return nil, err
	}

	p := &SQL{
		name:   name,
		db:     conn,
		alerts: &PostgresAlert{db: conn, tableCfg: alertsCfg, timeout: timeout, logger: logger},
		kv:     &PostgresKV{db: conn, tableCfg: kvCfg, timeout: timeout, logger: logger},
	}

	return p, nil
}

// Name returns the storage name
func (p *SQL) Name() string {
	return p.name
}

// Stop the storage
func (p *SQL) Stop() error {
	return p.db.Close()
}

// KV returns KV storage
func (p *SQL) KV() corestorage.KV {
	return p.kv
}

// Alert returns Alert storage
func (p *SQL) Alert() corestorage.Alert {
	return p.alerts
}
