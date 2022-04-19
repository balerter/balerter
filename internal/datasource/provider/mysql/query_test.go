package mysql

import (
	"math/rand"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/luaformatter"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

var _db *sqlx.DB

func getDB(t *testing.T) *sqlx.DB {
	if _db != nil {
		return _db
	}
	rand.Seed(time.Now().UnixNano())

	var err error
	_db, err = sqlx.Connect("mysql", "mysql:secret@tcp(127.0.0.1:3306)/db")
	require.NoError(t, err)
	return _db
}

func TestQuery_ErrorQuery(t *testing.T) {
	m := &MySQL{
		logger:  zap.NewNop(),
		timeout: time.Second,
	}

	query := "simple query"

	db := getDB(t)

	m.db = db

	luaState := lua.NewState()
	luaState.Push(lua.LString(query))

	n := m.query(luaState)

	assert.Equal(t, 2, n)

	e := luaState.Get(3)
	assert.Equal(t, lua.LTString, e.Type())
	assert.Equal(t, "Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'simple query' at line 1", e.String())
}

func TestQuery(t *testing.T) {
	db := getDB(t)

	query := "SELECT '1' AS ID,'Foo' AS Name UNION SELECT '20','Bar'"

	m := &MySQL{
		logger:  zap.NewNop(),
		timeout: time.Second,
		db:      db,
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
