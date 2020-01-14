package clickhouse

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/balerter/balerter/internal/config"
	"github.com/jmoiron/sqlx"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"time"
)

type Clickhouse struct {
	name string
	db   *sqlx.DB
}

func New(cfg config.DataSourceClickhouse) (*Clickhouse, error) {
	c := &Clickhouse{
		name: cfg.Name,
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

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second)
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

func (m *Clickhouse) GetLoader() lua.LGFunction {
	return m.loader
}

func (m *Clickhouse) loader(L *lua.LState) int {
	return 0
}
