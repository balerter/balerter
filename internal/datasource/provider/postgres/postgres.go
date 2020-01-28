package postgres

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type Postgres struct {
	name   string
	logger *zap.Logger
	db     *sqlx.DB
}

func New(cfg config.DataSourcePostgres, logger *zap.Logger) (*Postgres, error) {
	p := &Postgres{
		name:   "postgres." + cfg.Name,
		logger: logger,
	}

	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&sslrootcert=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
		cfg.SSLCertPath,
	)
	var err error

	p.db, err = sqlx.Connect("postgres", pgConnString)
	if err != nil {
		return nil, err
	}

	if err := p.db.Ping(); err != nil {
		p.db.Close()
		return nil, err
	}

	return p, nil
}

func (m *Postgres) Stop() error {
	return m.db.Close()
}

func (m *Postgres) Name() string {
	return m.name
}

func (m *Postgres) GetLoader(_ *script.Script) lua.LGFunction {
	return m.loader
}

func (m *Postgres) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.query,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}
