package postgres

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/postgres"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Postgres implements CoreStorage with the Postgres as a storage backend
type Postgres struct {
	name   string
	db     *sqlx.DB
	alerts *PostgresAlert
	kv     *PostgresKV
}

func New(cfg *postgres.Postgres, logger *zap.Logger) (*Postgres, error) {
	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&sslrootcert=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
		cfg.SSLCertPath,
	)
	conn, err := sqlx.Connect("postgres", pgConnString)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	p := &Postgres{
		name:   "postgres." + cfg.Name,
		db:     conn,
		alerts: &PostgresAlert{db: conn, table: cfg.TableAlerts, timeout: cfg.Timeout, logger: logger},
		kv:     &PostgresKV{db: conn, table: cfg.TableKV, timeout: cfg.Timeout, logger: logger},
	}

	return p, nil
}

func (p *Postgres) Name() string {
	return p.name
}

func (p *Postgres) Stop() error {
	return p.db.Close()
}

func (p *Postgres) KV() corestorage.KV {
	return p.kv
}

func (p *Postgres) Alert() corestorage.Alert {
	return p.alerts
}
