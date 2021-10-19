package postgres

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestNew_ErrorConnect(t *testing.T) {
	connFunc := func(_ context.Context, _ string) (*pgxpool.Pool, error) {
		return nil, fmt.Errorf("err1")
	}

	cfg := postgres.Postgres{}

	_, err := New(cfg, connFunc, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "error connect to to postgres, err1", err.Error())
}

func TestNew(t *testing.T) {
	cfg := postgres.Postgres{
		Name:        "pg1",
		Host:        "127.0.0.1",
		Port:        35432,
		Username:    "postgres",
		Password:    "secret",
		Database:    "db",
		SSLMode:     "disable",
		SSLCertPath: "",
		Timeout:     10,
	}

	p, err := New(cfg, pgxpool.Connect, zap.NewNop())

	require.NoError(t, err)
	assert.IsType(t, &Postgres{}, p)
}

func TestName(t *testing.T) {
	p := &Postgres{name: "Foo"}
	assert.Equal(t, "Foo", p.Name())
}

func TestStop(t *testing.T) {
	dbmock := &dbpoolMock{
		CloseFunc: func() {},
	}

	p := &Postgres{
		db: dbmock,
	}

	err := p.Stop()
	require.NoError(t, err)

	assert.Equal(t, 1, len(dbmock.CloseCalls()))
}

func TestGetLoader(t *testing.T) {
	p := &Postgres{}

	loader := p.GetLoader(nil)

	luaState := lua.NewState()

	n := loader(luaState)
	assert.Equal(t, 1, n)

	v := luaState.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func TestModuleName(t *testing.T) {
	assert.Equal(t, "postgres.Foo", ModuleName("Foo"))
}
