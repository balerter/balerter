package mysql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/luaformatter"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestQuery_ErrorQuery(t *testing.T) {
	m := &MySQL{
		logger:  zap.NewNop(),
		timeout: time.Second,
	}
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	query := "simple query"

	mock.ExpectQuery(query).WillReturnError(fmt.Errorf("err1"))

	m.db = sqlx.NewDb(mockDB, "sqlmock")

	luaState := lua.NewState()
	luaState.Push(lua.LString(query))

	n := m.query(luaState)

	assert.Equal(t, 2, n)

	e := luaState.Get(3)
	assert.Equal(t, lua.LTString, e.Type())
	assert.Equal(t, "err1", e.String())
}

func TestQuery(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := "simple query"

	rows := sqlmock.NewRows([]string{"ID", "Name"}).
		AddRow(1, "Foo").
		AddRow(20, "Bar")

	dbmock.ExpectQuery(query).WillReturnRows(rows)

	m := &MySQL{
		logger:  zap.NewNop(),
		timeout: time.Second,
		db:      sqlx.NewDb(db, "sqlmock"),
	}

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

	assert.Equal(t, "{\"ID\":\"1\",\"Name\":\"Foo\"}", row1str)
	assert.Equal(t, "{\"ID\":\"20\",\"Name\":\"Bar\"}", row2str)
}

func TestQuery_RowNextError(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := "simple query"

	rows := sqlmock.NewRows([]string{"ID", "Name"}).
		AddRow(1, "Foo").
		AddRow(2, "Bar").
		RowError(1, fmt.Errorf("err2"))

	dbmock.ExpectQuery(query).WillReturnRows(rows)

	m := &MySQL{
		logger:  zap.NewNop(),
		timeout: time.Second,
		db:      sqlx.NewDb(db, "sqlmock"),
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString(query))

	n := m.query(luaState)

	assert.Equal(t, 2, n)

	arg2 := luaState.Get(2)
	arg3 := luaState.Get(3)

	assert.Equal(t, arg2.Type(), lua.LTNil)
	assert.Equal(t, arg3.Type(), lua.LTString)
	assert.Equal(t, "error next: err2", arg3.String())
}
