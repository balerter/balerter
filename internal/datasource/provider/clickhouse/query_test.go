package clickhouse

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
	"time"
)

var _db *sqlx.DB

func getDB(t *testing.T) *sqlx.DB {
	if _db != nil {
		return _db
	}

	connString := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s&%s",
		"127.0.0.1",
		9000,
		"default",
		"",
		"default",
		"",
	)

	var err error

	_db, err = sqlx.ConnectContext(context.Background(), "clickhouse", connString)
	require.NoError(t, err)
	return _db
}

func TestClickhouse_query(t *testing.T) {
	ch := &Clickhouse{
		db:      getDB(t),
		logger:  zap.NewNop(),
		timeout: time.Second,
	}

	state := lua.NewState()
	state.Push(lua.LString("SELECT 1+1 AS num"))

	n := ch.query(state)
	assert.Equal(t, 2, n)

	v := state.Get(2)
	require.Equal(t, lua.LTTable, v.Type())

	vv := v.(*lua.LTable)

	var found bool

	vv.ForEach(func(value lua.LValue, value2 lua.LValue) {
		tbl := value2.(*lua.LTable)
		num := tbl.RawGetString("num")
		found = "2" == num.String()
	})
	assert.True(t, found)
}
