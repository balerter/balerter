package sql

import (
	"fmt"
	"strings"
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
		err2 := conn.Close()
		if err2 != nil {
			return nil, fmt.Errorf("error close sql connection after wrong ping %v, %w", err2, err)
		}
		return nil, err
	}

	p := &SQL{
		name:   name,
		db:     conn,
		alerts: &PostgresAlert{db: conn, table: tableAlerts, timeout: timeout, logger: logger},
		kv:     &PostgresKV{db: conn, table: tableKV, timeout: timeout, logger: logger},
	}

	err = p.createTableAlerts(tableAlerts)
	if err != nil {
		return nil, fmt.Errorf("error create alerts table, %w", err)
	}

	err = p.createTableKV(tableKV)
	if err != nil {
		return nil, fmt.Errorf("error create kv table, %w", err)
	}

	return p, nil
}

func (p *SQL) createTableKV(table string) error {
	query := `create table if not exists {%TABLE%}
(
	key varchar not null primary key,
	value text
);
`

	query = strings.Replace(query, "{%TABLE%}", table, -1)
	_, err := p.db.Exec(query)
	return err
}

func (p *SQL) createTableAlerts(table string) error {
	query := `
create table if not exists {%TABLE%}
(
	id varchar not null primary key,
	level int4 default 0 not null,
	count int4 default 0,
	last_change timestamp default CURRENT_TIMESTAMP,
	start timestamp default CURRENT_TIMESTAMP
);
`

	query = strings.Replace(query, "{%TABLE%}", table, -1)

	_, err := p.db.Exec(query)
	return err
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
