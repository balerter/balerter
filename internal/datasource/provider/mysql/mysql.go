package mysql

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	_ "github.com/go-sql-driver/mysql" // import DB driver
	"github.com/jmoiron/sqlx"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

func ModuleName(name string) string {
	return "mysql." + name
}

func Methods() []string {
	return []string{
		"query",
	}
}

type MySQL struct {
	name    string
	logger  *zap.Logger
	db      *sqlx.DB
	timeout time.Duration
}

type SQLConnFunc func(string, string) (*sqlx.DB, error)

func New(cfg *config.DataSourceMysql, sqlConnFunc SQLConnFunc, logger *zap.Logger) (*MySQL, error) {
	p := &MySQL{
		name:    ModuleName(cfg.Name),
		logger:  logger,
		timeout: cfg.Timeout,
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

func (m *MySQL) Stop() error {
	return m.db.Close()
}

func (m *MySQL) Name() string {
	return m.name
}

func (m *MySQL) GetLoader(_ *script.Script) lua.LGFunction {
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
