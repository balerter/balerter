package postgres

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/luaformatter"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestQuery_ErrorQuery(t *testing.T) {
	query := "simple query"

	dbmock := &dbpoolMock{
		QueryFunc: func(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
			require.Equal(t, query, sql)
			return nil, fmt.Errorf("err1")
		},
	}

	m := &Postgres{
		logger:  zap.NewNop(),
		timeout: time.Second,
		db:      dbmock,
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString(query))

	n := m.query(luaState)

	assert.Equal(t, 2, n)

	e := luaState.Get(3)
	assert.Equal(t, lua.LTString, e.Type())
	assert.Equal(t, "err1", e.String())
}

func TestQuery(t *testing.T) {
	cfg := postgres.Postgres{
		Name:        "pg1",
		Host:        "127.0.0.1",
		Port:        35432,
		Username:    "postgres",
		Password:    "secret",
		Database:    "db",
		SSLMode:     "disable",
		SSLCertPath: "",
		Timeout:     10000,
	}

	m, err := New(cfg, pgxpool.Connect, zap.NewNop())
	require.NoError(t, err)

	query := "select * from (values (1, 'Foo', true, null), (20, 'Bar', false, TIMESTAMP '2004-10-19 10:23:54+02')) as t(id, name, is_male, birthday)"

	luaState := lua.NewState()
	luaState.Push(lua.LString(query))

	n := m.query(luaState)

	assert.Equal(t, 2, n)

	arg2 := luaState.Get(2)
	arg3 := luaState.Get(3)

	assert.Equal(t, arg3.Type(), lua.LTNil)
	assert.Equal(t, arg2.Type(), lua.LTTable)

	n = arg2.(*lua.LTable).Len()
	assert.Equal(t, 2, n)
	row1 := arg2.(*lua.LTable).RawGet(lua.LNumber(1))
	row2 := arg2.(*lua.LTable).RawGet(lua.LNumber(2))

	row1str, err := luaformatter.TableToString(row1.(*lua.LTable))
	require.NoError(t, err)
	row2str, err := luaformatter.TableToString(row2.(*lua.LTable))
	require.NoError(t, err)

	assert.Equal(t, `{"birthday":"<nil>","id":"1","is_male":"true","name":"Foo"}`, row1str)
	assert.Equal(t, `{"birthday":"2004-10-19 10:23:54 +0000 UTC","id":"20","is_male":"false","name":"Bar"}`, row2str)
}
