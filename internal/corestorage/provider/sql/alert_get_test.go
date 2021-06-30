package sql

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/alert"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGet_error_query_row(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE id = \$1`).WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, err = pa.Get("foo")
	require.Error(t, err)
	assert.Equal(t, "error select alert, err1", err.Error())
}

func TestGet_error_row_scan(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "level", "count", "last_change", "start"}).
		AddRow(
			sql.NullString{String: "not valid value", Valid: false},
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE id = \$1`).WillReturnRows(rows)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, err = pa.Get("foo")
	require.Error(t, err)
	assert.Equal(t, "error scan result, sql: Scan error on column index 0, name \"id\": converting NULL to string is unsupported", err.Error())
}

func TestGet_no_rows(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "level", "count", "last_change", "start"})

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE id = \$1`).WillReturnRows(rows)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	a, err := pa.Get("foo")
	assert.Nil(t, a)
	assert.Nil(t, err)
}

func TestGet_error_parse_level(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE id = \$1`).WillReturnRows(rows)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, err = pa.Get("foo")
	require.Error(t, err)
	assert.Equal(t, "error parse level 10 for alert 10, bad level", err.Error())
}

func TestGet(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: now, Valid: true},
			sql.NullTime{Time: now.Add(time.Hour), Valid: true},
		)

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE id = \$1`).WillReturnRows(rows)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	a, err := pa.Get("foo")
	require.NoError(t, err)
	require.NotNil(t, a)

	assert.Equal(t, "10", a.Name)
	assert.Equal(t, alert.LevelSuccess, a.Level)
	assert.Equal(t, 1, a.Count)
	assert.Equal(t, now.Unix(), a.LastChange.Unix())
	assert.Equal(t, now.Add(time.Hour).Unix(), a.Start.Unix())
}
