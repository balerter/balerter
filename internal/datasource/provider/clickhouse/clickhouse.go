package clickhouse

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	clickhouseCfg "github.com/balerter/balerter/internal/config/datasources/clickhouse"
	"github.com/balerter/balerter/internal/modules"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/jmoiron/sqlx"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

// ModuleName returns the module name
func ModuleName(name string) string {
	return "clickhouse." + name
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"query",
	}
}

type dbConnection interface {
	Ping() error
	Close() error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Clickhouse represents datasource of type Clickhouse
type Clickhouse struct {
	name    string
	logger  *zap.Logger
	db      dbConnection
	timeout time.Duration
}

// New creates new Clickhouse datasource
func New(cfg clickhouseCfg.Clickhouse, logger *zap.Logger) (*Clickhouse, error) {
	c := &Clickhouse{
		name:    ModuleName(cfg.Name),
		logger:  logger,
		timeout: time.Millisecond * time.Duration(cfg.Timeout),
	}

	if c.timeout == 0 {
		c.timeout = defaultTimeout
	}

	chSecureString := "secure=false"

	if cfg.SSLCertPath != "" {
		caCertPool := x509.NewCertPool()
		caCert, err := os.ReadFile(cfg.SSLCertPath)
		if err != nil {
			return nil, fmt.Errorf("error load clickhouse cert file, %v", err)
		}

		caCertPool.AppendCertsFromPEM(caCert)

		if err := clickhouse.RegisterTLSConfig("chtls", &tls.Config{
			RootCAs: caCertPool,
		}); err != nil {
			return nil, fmt.Errorf("error register tls config, %v", err)
		}

		chSecureString = "secure=true&tls_config=chtls"
	}

	connString := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s&%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		chSecureString,
	)

	ctx, ctxCancel := context.WithTimeout(context.Background(), c.timeout)
	defer ctxCancel()

	var err error

	if c.db, err = sqlx.ConnectContext(ctx, "clickhouse", connString); err != nil {
		return nil, fmt.Errorf("error connect to clickhouse, %v", err)
	}

	if err := c.db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping clickhouse, %v", err)
	}

	return c, nil
}

// Stop the datasource
func (m *Clickhouse) Stop() error {
	return m.db.Close()
}

// Name returns the datasource name
func (m *Clickhouse) Name() string {
	return m.name
}

func (m *Clickhouse) GetLoaderJS(_ modules.Job) require.ModuleLoader {
	return func(runtime *goja.Runtime, object *goja.Object) {

	}
}

// GetLoader returns the datasource lua loader
func (m *Clickhouse) GetLoader(_ modules.Job) lua.LGFunction {
	return m.loader
}

func (m *Clickhouse) loader(luaState *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.query,
	}

	mod := luaState.SetFuncs(luaState.NewTable(), exports)

	luaState.Push(mod)
	return 1
}
