package postgres

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // DB driver
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

func ModuleName(name string) string {
	return "postgres." + name
}

func Methods() []string {
	return []string{
		"query",
	}
}

type Postgres struct {
	name    string
	logger  *zap.Logger
	db      *sqlx.DB
	timeout time.Duration
}

type SQLConnFunc func(string, string) (*sqlx.DB, error)

func New(cfg postgres.Postgres, sqlConnFunc SQLConnFunc, logger *zap.Logger) (*Postgres, error) {
	p := &Postgres{
		name:    ModuleName(cfg.Name),
		logger:  logger,
		timeout: time.Millisecond * time.Duration(cfg.Timeout),
	}

	if p.timeout == 0 {
		p.timeout = defaultTimeout
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

	p.db, err = sqlConnFunc("postgres", pgConnString)
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

func (m *Postgres) loader(luaState *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.query,
	}

	mod := luaState.SetFuncs(luaState.NewTable(), exports)

	luaState.Push(mod)
	return 1
}
