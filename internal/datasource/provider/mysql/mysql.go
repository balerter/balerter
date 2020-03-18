package mysql

import (
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

type MySQL struct {
	name    string
	logger  *zap.Logger
	db      *sqlx.DB
	timeout time.Duration
}

func New(cfg config.DataSourceMysql, logger *zap.Logger) (*MySQL, error) {
	p := &MySQL{
		name:    "mysql." + cfg.Name,
		logger:  logger,
		timeout: cfg.Timeout,
	}

	if p.timeout == 0 {
		p.timeout = defaultTimeout
	}

	var err error

	p.db, err = sqlx.Connect("mysql", cfg.DSN)
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

func (m *MySQL) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.query,
	}

	mod := L.SetFuncs(L.NewTable(), exports)

	L.Push(mod)
	return 1
}
