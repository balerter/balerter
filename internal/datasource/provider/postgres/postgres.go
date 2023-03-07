package postgres

//go:generate moq -out dbpool_mock.go -skip-ensure -fmt goimports . dbpool

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"time"

	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/modules"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
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

type dbpool interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Close()
}

// Postgres represents the datasource of the type Postgres
type Postgres struct {
	name    string
	logger  *zap.Logger
	db      dbpool
	timeout time.Duration
}

type SQLConnectFunc func(ctx context.Context, connString string) (*pgxpool.Pool, error)

// New creates new Postgres datasource
func New(cfg postgres.Postgres, connFunc SQLConnectFunc, logger *zap.Logger) (*Postgres, error) {
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

	db, errConnect := connFunc(context.Background(), pgConnString)
	if errConnect != nil {
		return nil, fmt.Errorf("error connect to to postgres, %w", errConnect)
	}

	p.db = db

	return p, nil
}

// Stop the datasource
func (m *Postgres) Stop() error {
	m.db.Close()
	return nil
}

// Name returns the datasource name
func (m *Postgres) Name() string {
	return m.name
}

func (m *Postgres) GetLoaderJS(_ modules.Job) require.ModuleLoader {
	return func(runtime *goja.Runtime, object *goja.Object) {

	}
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
