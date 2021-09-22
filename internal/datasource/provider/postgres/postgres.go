package postgres

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/modules"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // DB driver
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

// ModuleName returns the module name
func ModuleName(name string) string {
	return "postgres." + name
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"query",
	}
}

// Postgres represents the datasource of the type Postgres
type Postgres struct {
	name    string
	logger  *zap.Logger
	db      *sqlx.DB
	timeout time.Duration
}

// SQLConnFunc represent SQL connection func
type SQLConnFunc func(string, string) (*sqlx.DB, error)

// New creates new Postgres datasource
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

// Stop the datasource
func (m *Postgres) Stop() error {
	return m.db.Close()
}

// Name returns the datasource name
func (m *Postgres) Name() string {
	return m.name
}

// GetLoader returns the datasource lua loader
func (m *Postgres) GetLoader(_ modules.Job) lua.LGFunction {
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
