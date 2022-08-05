package mysql

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/datasources/mysql"
	"github.com/balerter/balerter/internal/modules"
	_ "github.com/go-sql-driver/mysql" // import DB driver
	"github.com/jmoiron/sqlx"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"net/http"
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

func (m *MySQL) CoreApiHandler(req []string, body []byte) (any, int, error) {
	return nil, http.StatusNotImplemented, fmt.Errorf("not implemented")
}

// Stop the datasource
func (m *MySQL) Stop() error {
	return m.db.Close()
}

// Name returns the datasource name
func (m *MySQL) Name() string {
	return m.name
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
