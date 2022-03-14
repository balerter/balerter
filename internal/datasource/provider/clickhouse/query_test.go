package clickhouse

import (
	"context"
	"fmt"
	"testing"
	"time"

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

	connString := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s&%s",
		"127.0.0.1",
		9000,
		"default",
		"",
		"default",
		"",
	)

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ctxCancel()

	var err error

	_db, err = sqlx.ConnectContext(ctx, "clickhouse", connString)
	require.NoError(t, err)
	return _db
}

func TestQuery(t *testing.T) {
	query := "SELECT 42 AS age, 'Foo' AS name UNION ALL SELECT 12, 'Bar'"

	db := getDB(t)

	ch := &Clickhouse{
		db:      db,
		logger:  zap.NewNop(),
		timeout: time.Second,
	}

	L := lua.NewState()
	L.Push(lua.LString(query))

	n := ch.query(L)
	assert.Equal(t, 2, n)

	arg2 := L.Get(2)
	arg3 := L.Get(3)

	assert.Equal(t, arg3.Type(), lua.LTNil)
	assert.Equal(t, arg2.Type(), lua.LTTable)

	type resultItem struct {
		age  string
		name string
	}

	results := []resultItem{
		{"12", "Bar"},
		{"42", "Foo"},
	}

	arg2.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
		require.Equal(t, value2.Type(), lua.LTTable)
		item := results[0]
		results = results[1:]

		value2.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
			key := value.String()
			v := value2.String()
			require.Contains(t, []string{"age", "name"}, key)
			switch key {
			case "age":
				assert.Equal(t, item.age, v)
			case "name":
				assert.Equal(t, item.name, v)
			}
		})
	})
}
