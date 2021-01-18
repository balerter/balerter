package postgres

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

func TestAlert_Update_begin_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin().WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelSuccess)
	require.Error(t, err)
	assert.Equal(t, "error start tx, err1", err.Error())
}

func TestAlert_Update_insert_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelSuccess)
	require.Error(t, err)
	assert.Equal(t, "error insert row, err1", err.Error())
}

func TestAlert_Update_affected_rows_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewErrorResult(fmt.Errorf("err1"))

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelSuccess)
	require.Error(t, err)
	assert.Equal(t, "error get affected rows count, err1", err.Error())
}

func TestAlert_Update_affected_1_error_commit(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 1)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectCommit().WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelSuccess)
	require.Error(t, err)
	assert.Equal(t, "error commit tx, err1", err.Error())
}

func TestAlert_Update_affected_1__not_change(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 1)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	a, updated, err := pa.Update("foo", alert.LevelSuccess)
	require.NoError(t, err)
	assert.False(t, updated)
	assert.Equal(t, 0, a.Count)
	assert.Equal(t, "foo", a.Name)
	assert.Equal(t, alert.LevelSuccess, a.Level)
}

func TestAlert_Update_affected_1__change(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 1)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	a, updated, err := pa.Update("foo", alert.LevelWarn)
	require.NoError(t, err)
	assert.True(t, updated)
	assert.Equal(t, 0, a.Count)
	assert.Equal(t, "foo", a.Name)
	assert.Equal(t, alert.LevelWarn, a.Level)
}

func TestAlert_Update_error_selectscan(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "error scan row, err1", err.Error())
}

func TestAlert_Update_error_level_from_int_with_commit_err(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"id", "level", "last_change", "start"}).AddRow(sql.NullInt64{Int64: 10, Valid: true}, sql.NullInt64{Int64: 10, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true})

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)
	sqlm.ExpectCommit().WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "error convert level, error commit tx, err1", err.Error())
}

func TestAlert_Update_error_level_from_int(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"id", "level", "last_change", "start"}).AddRow(sql.NullInt64{Int64: 10, Valid: true}, sql.NullInt64{Int64: 10, Valid: true}, sql.NullTime{Time: time.Now(), Valid: true}, sql.NullTime{Time: time.Now(), Valid: true})

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES \(\$1, \$2, 1, NOW\(\), NOW\(\)\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		table:  "alerts",
		db:     dbx,
		logger: zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "error convert level, bad level", err.Error())
}
