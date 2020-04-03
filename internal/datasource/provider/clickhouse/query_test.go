package clickhouse

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := "some query"

	rows := sqlmock.NewRows([]string{"age", "name"}).
		AddRow(42, "Foo").
		AddRow(12, "Bar")

	dbmock.ExpectQuery(query).WillReturnRows(rows)

	ch := &Clickhouse{
		db:      sqlx.NewDb(db, "sqlmock"),
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
		{"42", "Foo"},
		{"12", "Bar"},
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
