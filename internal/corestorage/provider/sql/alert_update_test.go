package sql

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
	"time"
)

func createTableAlertsCfg() tables.TableAlerts {
	return tables.TableAlerts{
		Table: "alerts",
		Fields: tables.AlertFields{
			Name:      "id",
			Level:     "level",
			Count:     "count",
			UpdatedAt: "last_change",
			CreatedAt: "start",
		},
	}
}

func TestAlert_Update_begin_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin().WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
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
	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).
		WillReturnError(fmt.Errorf("err1"))
	sqlm.ExpectRollback().WillReturnError(fmt.Errorf("err2"))

	core, logger := observer.New(zap.DebugLevel)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.New(core),
	}

	_, _, err = pa.Update("foo", alert.LevelSuccess)
	require.Error(t, err)
	assert.Equal(t, "error insert row, err1", err.Error())
	assert.Equal(t, 1, logger.FilterMessage("error rollback tx").Len())
}

func TestAlert_Update_affected_rows_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewErrorResult(fmt.Errorf("err1"))

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
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

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectCommit().WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
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

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
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

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
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

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
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

	rows := sqlmock.NewRows([]string{"id", "level", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)
	sqlm.ExpectCommit().WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "error convert level 10, bad level", err.Error())
}

func TestAlert_Update_error_level_from_int(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullInt64{Int64: 10, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "error convert level 10, bad level", err.Error())
}

func TestAlert_Update_error_update_on_not_change_level(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)

	sqlm.ExpectExec(`UPDATE alerts SET count = count \+ 1, last_change = CURRENT_TIMESTAMP WHERE id = \$1`).WillReturnError(fmt.Errorf("err1"))
	sqlm.ExpectRollback()

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "error update row, err1", err.Error())
}

func TestAlert_Update_error_update_on_not_change_level_with_rollback_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)

	sqlm.ExpectExec(`UPDATE alerts SET count = count \+ 1, ` +
		`last_change = CURRENT_TIMESTAMP WHERE id = \$1`).WillReturnError(fmt.Errorf("err1"))
	sqlm.ExpectRollback().WillReturnError(fmt.Errorf("err1"))

	core, logger := observer.New(zap.DebugLevel)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.New(core),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "error update row, err1", err.Error())
	assert.Equal(t, 1, logger.FilterMessage("error rollback tx").Len())
}

func TestAlert_Update_not_change_level_error_commit(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)

	sqlm.ExpectExec(`UPDATE alerts SET count = count \+ 1, ` +
		`last_change = CURRENT_TIMESTAMP WHERE id = \$1`).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlm.ExpectCommit().WillReturnError(fmt.Errorf("err1"))

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelWarn)
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestAlert_Update_not_change_level(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)

	sqlm.ExpectExec(`UPDATE alerts SET count = count \+ 1, ` +
		`last_change = CURRENT_TIMESTAMP WHERE id = \$1`).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	alertName := "foo22"

	a, updated, err := pa.Update(alertName, alert.LevelWarn)
	require.NoError(t, err)
	assert.False(t, updated)

	assert.Equal(t, alertName, a.Name)
	assert.Equal(t, alert.LevelWarn, a.Level)
	assert.Equal(t, 2, a.Count)

	v, err := metrics.GetAlertLevel(alertName)
	require.NoError(t, err)
	assert.Equal(t, float64(alert.LevelWarn), v)
}

func TestAlert_Update_change_level_exec_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)

	sqlm.ExpectExec(`UPDATE alerts SET level = \$1, count = 1, ` +
		`last_change = CURRENT_TIMESTAMP WHERE id = \$2`).WillReturnError(fmt.Errorf("err1"))
	sqlm.ExpectRollback()

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	_, _, err = pa.Update("foo", alert.LevelError)
	require.Error(t, err)
	assert.Equal(t, "error update row, err1", err.Error())
}

func TestAlert_Update_change_level_exec_error_with_error_rollback(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)

	sqlm.ExpectExec(`UPDATE alerts SET level = \$1, count = 1, ` +
		`last_change = CURRENT_TIMESTAMP WHERE id = \$2`).WillReturnError(fmt.Errorf("err1"))
	sqlm.ExpectRollback().WillReturnError(fmt.Errorf("err2"))

	core, logger := observer.New(zap.DebugLevel)

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.New(core),
	}

	_, _, err = pa.Update("foo", alert.LevelError)
	require.Error(t, err)
	assert.Equal(t, "error update row, err1", err.Error())
	assert.Equal(t, 1, logger.FilterMessage("error rollback tx").Len())
}

func TestAlert_Update_change_level(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectBegin()

	r := sqlmock.NewResult(1, 0)

	rows := sqlmock.NewRows([]string{"level", "count", "last_change", "start"}).
		AddRow(
			sql.NullInt64{Int64: 2, Valid: true},
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
			sql.NullTime{Time: time.Now(), Valid: true},
		)

	sqlm.ExpectExec(`INSERT INTO alerts \(id, level, count, last_change, start\) VALUES ` +
		`\(\$1, \$2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) ON CONFLICT \(id\) DO NOTHING`).WillReturnResult(r)
	sqlm.ExpectQuery(`SELECT`).WillReturnRows(rows)

	sqlm.ExpectExec(`UPDATE alerts SET level = \$1, count = 1, ` +
		`last_change = CURRENT_TIMESTAMP WHERE id = \$2`).WillReturnResult(sqlmock.NewResult(1, 1))
	sqlm.ExpectCommit()

	pa := &PostgresAlert{
		tableCfg: createTableAlertsCfg(),
		db:       dbx,
		logger:   zap.NewNop(),
	}

	alertName := "foo33"

	a, updated, err := pa.Update(alertName, alert.LevelError)
	require.NoError(t, err)
	assert.True(t, updated)

	assert.Equal(t, alertName, a.Name)
	assert.Equal(t, alert.LevelError, a.Level)
	assert.Equal(t, 0, a.Count)
}
