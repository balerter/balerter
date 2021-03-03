package mysql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/config/datasources/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestNew_ErrorConnect(t *testing.T) {
	mockConnFunc := func(string, string) (*sqlx.DB, error) {
		return nil, fmt.Errorf("err1")
	}

	cfg := mysql.Mysql{}

	_, err := New(cfg, mockConnFunc, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestNew_ErrorPing(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	mockConnFunc := func(string, string) (*sqlx.DB, error) {
		return sqlx.NewDb(db, "sqlmock"), nil
	}

	dbmock.ExpectPing().WillReturnError(fmt.Errorf("err2"))

	cfg := mysql.Mysql{}

	_, err = New(cfg, mockConnFunc, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "err2", err.Error())
}

func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)

	mockConnFunc := func(string, string) (*sqlx.DB, error) {
		return sqlx.NewDb(db, "sqlmock"), nil
	}

	cfg := mysql.Mysql{}

	p, err := New(cfg, mockConnFunc, zap.NewNop())

	require.NoError(t, err)
	assert.IsType(t, &MySQL{}, p)
}

func TestName(t *testing.T) {
	p := &MySQL{name: "Foo"}
	assert.Equal(t, "Foo", p.Name())
}

func TestStop_Error(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	p := &MySQL{
		db: sqlx.NewDb(db, "sqlmock"),
	}

	dbmock.ExpectClose().WillReturnError(fmt.Errorf("err1"))

	err = p.Stop()
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestStop(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	p := &MySQL{
		db: sqlx.NewDb(db, "sqlmock"),
	}

	dbmock.ExpectClose()

	err = p.Stop()
	require.NoError(t, err)
}

func TestGetLoader(t *testing.T) {
	p := &MySQL{}

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
	assert.Equal(t, "mysql.Foo", ModuleName("Foo"))
}
