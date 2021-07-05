package clickhouse

import (
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

func TestStop(t *testing.T) {
	cc := &dbConnectionMock{
		CloseFunc: func() error {
			return fmt.Errorf("err1")
		},
	}
	ch := &Clickhouse{
		db: cc,
	}

	err := ch.Stop()

	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())

	assert.Equal(t, 1, len(cc.CloseCalls()))
}

func TestNew_error_connect(t *testing.T) {
	_, err := New(clickhouseCfg.Clickhouse{}, zap.NewNop())
	require.Error(t, err)
	assert.Equal(t, "error connect to clickhouse, dial tcp :0: connect: can't assign requested address", err.Error())
}

func TestNew(t *testing.T) {
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
