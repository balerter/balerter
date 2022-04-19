package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	clickhouseCfg "github.com/balerter/balerter/internal/config/datasources/clickhouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestName(t *testing.T) {
	ch := &Clickhouse{name: "foo"}
	assert.Equal(t, "foo", ch.Name())
}

func TestLoader(t *testing.T) {
	ch := &Clickhouse{}

	f := ch.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func TestModuleName(t *testing.T) {
	assert.Equal(t, "clickhouse.ch1", ModuleName("ch1"))
}

func TestMethods(t *testing.T) {
	assert.Equal(t, []string{"query"}, Methods())
}

type dbMock struct {
	stopped bool
}

func (db *dbMock) Ping() error { return nil }
func (db *dbMock) Close() error {
	db.stopped = true
	return fmt.Errorf("err1")
}
func (db *dbMock) QueryContext(_ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func TestStop(t *testing.T) {
	db := &dbMock{}
	ch := &Clickhouse{
		db: db,
	}

	err := ch.Stop()

	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
	assert.True(t, db.stopped)
}

//func TestNew_error_connect(t *testing.T) {
//	_, err := New(clickhouseCfg.Clickhouse{}, zap.NewNop())
//	require.Error(t, err)
//	assert.Equal(t, "error connect to clickhouse, dial tcp :0: connect: can't assign requested address", err.Error())
//}

func TestNew(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ch, err := New(clickhouseCfg.Clickhouse{
		Name:        "ch1",
		Host:        "127.0.0.1",
		Port:        9000,
		Username:    "default",
		Password:    "",
		Database:    "default",
		SSLCertPath: "",
		Timeout:     0,
	}, zap.NewNop())

	require.NoError(t, err)
	assert.IsType(t, &Clickhouse{}, ch)
}

func TestNew_error_cert(t *testing.T) {
	_, err := New(clickhouseCfg.Clickhouse{
		Name:        "ch1",
		Host:        "127.0.0.1",
		Port:        9000,
		Username:    "default",
		Password:    "",
		Database:    "default",
		SSLCertPath: "notfound",
		Timeout:     0,
	}, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "error load clickhouse cert file, open notfound: no such file or directory", err.Error())
}
