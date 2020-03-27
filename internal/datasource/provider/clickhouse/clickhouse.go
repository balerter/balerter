package clickhouse

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/jmoiron/sqlx"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

var (
	defaultTimeout = time.Second * 5
)

func ModuleName(name string) string {
	return "clickhouse." + name
}

func Methods() []string {
	return []string{
		"query",
	}
}

type Clickhouse struct {
	name    string
	logger  *zap.Logger
	db      *sqlx.DB
	timeout time.Duration
}

func New(cfg config.DataSourceClickhouse, logger *zap.Logger) (*Clickhouse, error) {
	c := &Clickhouse{
		name:    ModuleName(cfg.Name),
		logger:  logger,
		timeout: cfg.Timeout,
	}

	if c.timeout == 0 {
		c.timeout = defaultTimeout
	}

	chSecureString := "secure=false"

	if cfg.SSLCertPath != "" {
		caCertPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile(cfg.SSLCertPath)
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

func (m *Clickhouse) Stop() error {
	return m.db.Close()
}

func (m *Clickhouse) Name() string {
	return m.name
}

func (m *Clickhouse) GetLoader(_ *script.Script) lua.LGFunction {
	return m.loader
}

func (m *Clickhouse) loader(L *lua.LState) int {
	var exports = map[string]lua.LGFunction{
		"query": m.query,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	// register other stuff
	//L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}
