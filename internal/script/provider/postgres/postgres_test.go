package postgres

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _db *sqlx.DB

func getDB(t *testing.T) *sqlx.DB {
	if _db != nil {
		return _db
	}
	rand.Seed(time.Now().UnixNano())

	var err error
	_db, err = sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&sslrootcert=%s",
		"postgres",
		"secret",
		"127.0.0.1",
		35432,
		"postgres",
		"disable",
		""))
	require.NoError(t, err)
	return _db
}

func TestProvider_Get_error_parse_meta(t *testing.T) {
	db := getDB(t)

	p := &Provider{
		query: "SELECT 'name', '-- @timeout abc\nfoo\nbar'",
		db:    db,
	}

	_, err := p.Get()
	require.Error(t, err)
	assert.Equal(t, "error parse 'abc' to time duration, time: invalid duration \"abc\"", err.Error())
}

func TestProvider_Get(t *testing.T) {
	db := getDB(t)

	p := &Provider{
		query: "select * from (values ('name', 'foo'), ('name','bar')) as t (name, body)",
		db:    db,
	}

	ss, err := p.Get()
	require.NoError(t, err)
	assert.Equal(t, 2, len(ss))

	for _, s := range ss {
		if string(s.Body) != "foo" && string(s.Body) != "bar" {
			t.Fatal("unexpected file body")
		}
	}
}
