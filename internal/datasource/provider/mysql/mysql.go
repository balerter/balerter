package mysql

import (
	"github.com/balerter/balerter/internal/config/datasources/mysql"
	"github.com/balerter/balerter/internal/modules"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	_ "github.com/go-sql-driver/mysql" // import DB driver
	"github.com/jmoiron/sqlx"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

// ModuleName returns the module name
func ModuleName(name string) string {
	return "mysql." + name
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"query",
	}
}

// MySQL represent the datasource of type MySQL
type MySQL struct {
	name    string
	logger  *zap.Logger
	db      *sqlx.DB
	timeout time.Duration
}

// SQLConnFunc represent ConnFunc
type SQLConnFunc func(string, string) (*sqlx.DB, error)

// New creates new MySQL datasource
func New(cfg mysql.Mysql, sqlConnFunc SQLConnFunc, logger *zap.Logger) (*MySQL, error) {
	p := &MySQL{
		name:    ModuleName(cfg.Name),
		logger:  logger,
		timeout: time.Millisecond * time.Duration(cfg.Timeout),
	}

	if p.timeout == 0 {
		p.timeout = defaultTimeout
	}

	var err error

	p.db, err = sqlConnFunc("mysql", cfg.DSN)
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
func (m *MySQL) Stop() error {
	return m.db.Close()
}

// Name returns the datasource name
func (m *MySQL) Name() string {
	return m.name
}

func (m *MySQL) GetLoaderJS(_ modules.Job) require.ModuleLoader {
	return func(runtime *goja.Runtime, object *goja.Object) {

	}
}

// GetLoader returns the datasource lua loader
func (m *MySQL) GetLoader(_ modules.Job) lua.LGFunction {
	return m.loader
}

func (m *MySQL) loader(luaState *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.query,
	}

	mod := luaState.SetFuncs(luaState.NewTable(), exports)

	luaState.Push(mod)
	return 1
}
