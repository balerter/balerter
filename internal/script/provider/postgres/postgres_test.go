package postgres

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/config/scripts/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew_error_connect(t *testing.T) {
	_, err := New(postgres.Postgres{})
	require.Error(t, err)
	assert.Equal(t, "dial tcp [::1]:0: connect: can't assign requested address", err.Error())
}

func TestProvider_Get_error_query(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	dbmock.ExpectQuery("QUERY").WillReturnError(fmt.Errorf("err1"))

	p := &Provider{
		query: "QUERY",
		db:    sqlx.NewDb(db, "sqlmock"),
	}

	_, err = p.Get()
	require.Error(t, err)
	assert.Equal(t, "error db query, err1", err.Error())
}

func TestProvider_Get_error_scan(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"name", "body"}).
		AddRow(sql.NullString{String: "bar", Valid: false}, "foo").
		AddRow(sql.NullString{String: "bar", Valid: true}, "bar")

	dbmock.ExpectQuery("QUERY").WillReturnRows(rows)

	p := &Provider{
		query: "QUERY",
		db:    sqlx.NewDb(db, "sqlmock"),
	}

	_, err = p.Get()
	require.Error(t, err)
	assert.Equal(t, "sql: Scan error on column index 0, name \"name\": converting NULL to string is unsupported", err.Error())
}

func TestProvider_Get_error_parse_meta(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"name", "body"}).
		AddRow(sql.NullString{String: "bar", Valid: true}, "-- @timeout abc\nfoo").
		AddRow(sql.NullString{String: "bar", Valid: true}, "bar")

	dbmock.ExpectQuery("QUERY").WillReturnRows(rows)

	p := &Provider{
		query: "QUERY",
		db:    sqlx.NewDb(db, "sqlmock"),
	}

	_, err = p.Get()
	require.Error(t, err)
	assert.Equal(t, "error parse 'abc' to time duration, time: invalid duration \"abc\"", err.Error())
}

func TestProvider_Get(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"name", "body"}).
		AddRow(sql.NullString{String: "foo", Valid: true}, "foo").
		AddRow(sql.NullString{String: "bar", Valid: true}, "bar")

	dbmock.ExpectQuery("QUERY").WillReturnRows(rows)

	p := &Provider{
		query: "QUERY",
		db:    sqlx.NewDb(db, "sqlmock"),
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

func TestProvider_Stop(t *testing.T) {
	db, e, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	e.ExpectClose()

	p := &Provider{
		db: sqlx.NewDb(db, "sqlmock"),
	}

	err = p.Stop()
	require.NoError(t, err)
}
