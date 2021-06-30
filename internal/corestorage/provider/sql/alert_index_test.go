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

func TestIndex_error_select_row(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE level IN \(1,2\)`).WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, err = pa.Index([]alert.Level{alert.LevelSuccess, alert.LevelWarn})
	require.Error(t, err)
	assert.Equal(t, "error select rows, err1", err.Error())
}

func TestIndex_error_scan(t *testing.T) {
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

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE level IN \(1,2\)`).WillReturnRows(rows)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, err = pa.Index([]alert.Level{alert.LevelSuccess, alert.LevelWarn})
	require.Error(t, err)
	assert.Equal(t, "error scan result, sql: Scan error on column index 0, name \"id\": converting NULL to string is unsupported", err.Error())
}

func TestIndex_convert_level(t *testing.T) {
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

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE level IN \(1,2\)`).WillReturnRows(rows)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, err = pa.Index([]alert.Level{alert.LevelSuccess, alert.LevelWarn})
	require.Error(t, err)
	assert.Equal(t, "error parse level 10 for alert 10, bad level", err.Error())
}

func TestIndex(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: now, Valid: true},
			sql.NullTime{Time: now.Add(time.Hour), Valid: true},
		).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: now, Valid: true},
			sql.NullTime{Time: now.Add(time.Hour), Valid: true},
		)

	sqlm.ExpectQuery(`SELECT id, level, count, last_change, start FROM alerts WHERE level IN \(1,2\)`).WillReturnRows(rows)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	a, err := pa.Index([]alert.Level{alert.LevelSuccess, alert.LevelWarn})
	require.NoError(t, err)
	require.Equal(t, 2, len(a))

	assert.Equal(t, "1", a[0].Name)
	assert.Equal(t, alert.LevelSuccess, a[0].Level)
	assert.Equal(t, 1, a[0].Count)
	assert.Equal(t, now.Unix(), a[0].LastChange.Unix())
	assert.Equal(t, now.Add(time.Hour).Unix(), a[0].Start.Unix())

	assert.Equal(t, "2", a[1].Name)
	assert.Equal(t, alert.LevelWarn, a[1].Level)
	assert.Equal(t, 1, a[1].Count)
	assert.Equal(t, now.Unix(), a[1].LastChange.Unix())
	assert.Equal(t, now.Add(time.Hour).Unix(), a[1].Start.Unix())

}
