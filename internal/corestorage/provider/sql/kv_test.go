package sql

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPostgresKV_All_error_query(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectQuery(`SELECT key, value FROM kv`).WillReturnError(fmt.Errorf("err1"))

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	_, err = kv.All()
	require.Error(t, err)
	assert.Equal(t, "error sql query, err1", err.Error())
}

func TestPostgresKV_All_scan_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"key", "value"}).
		AddRow(
			sql.NullInt64{Int64: 10, Valid: false},
			sql.NullInt64{Int64: 10, Valid: true},
		)
	sqlm.ExpectQuery(`SELECT key, value FROM kv`).WillReturnRows(rows)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	_, err = kv.All()
	require.Error(t, err)
	assert.Equal(t, "error scan result, sql: Scan error on column index 0, "+
		"name \"key\": converting NULL to string is unsupported", err.Error())
}

func TestPostgresKV_All(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"key", "value"}).
		AddRow(
			sql.NullString{String: "a1", Valid: true},
			sql.NullString{String: "a2", Valid: true},
		).
		AddRow(
			sql.NullString{String: "b1", Valid: true},
			sql.NullString{String: "b2", Valid: true},
		)
	sqlm.ExpectQuery(`SELECT key, value FROM kv`).WillReturnRows(rows)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	data, err := kv.All()
	require.NoError(t, err)
	require.Equal(t, 2, len(data))

	e, ok := data["a1"]
	require.True(t, ok)
	assert.Equal(t, "a2", e)

	e, ok = data["b1"]
	require.True(t, ok)
	assert.Equal(t, "b2", e)
}

func TestPostgresKV_Put_error_query(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectExec(`INSERT INTO kv \(key, value\) VALUES ` +
		`\(\$1, \$2\) ON CONFLICT \(key\) DO NOTHING`).WillReturnError(fmt.Errorf("err1"))

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	err = kv.Put("k", "v")
	require.Error(t, err)
	assert.Equal(t, "error sql query, err1", err.Error())
}

func TestPostgresKV_Put_0_affected(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	r := sqlmock.NewResult(0, 0)

	sqlm.ExpectExec(`INSERT INTO kv \(key, value\) VALUES ` +
		`\(\$1, \$2\) ON CONFLICT \(key\) DO NOTHING`).WillReturnResult(r)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	err = kv.Put("k", "v")
	require.Error(t, err)
	assert.Equal(t, "key already exists", err.Error())
}

func TestPostgresKV_Put(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	r := sqlmock.NewResult(0, 1)

	sqlm.ExpectExec(`INSERT INTO kv \(key, value\) VALUES` +
		` \(\$1, \$2\) ON CONFLICT \(key\) DO NOTHING`).WillReturnResult(r)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	err = kv.Put("k", "v")
	require.NoError(t, err)
}

func TestPostgresKV_Get_error_query(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectQuery(`SELECT value FROM kv WHERE key = \$1`).WillReturnError(fmt.Errorf("err1"))

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	_, err = kv.Get("k")
	require.Error(t, err)
	assert.Equal(t, "error sql query, err1", err.Error())
}

func TestPostgresKV_Get_error_no_rows(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"value"})

	sqlm.ExpectQuery(`SELECT value FROM kv WHERE key = \$1`).WillReturnRows(rows)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	_, err = kv.Get("k")
	require.Error(t, err)
	assert.Equal(t, ErrNoRow.Error(), err.Error())
}

func TestPostgresKV_Get_error_scan(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"value"}).
		AddRow(
			sql.NullString{String: "not valid value", Valid: false},
		)

	sqlm.ExpectQuery(`SELECT value FROM kv WHERE key = \$1`).WillReturnRows(rows)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	_, err = kv.Get("k")
	require.Error(t, err)
	assert.Equal(t, "error scan result, sql: Scan error on column index 0, "+
		"name \"value\": converting NULL to string is unsupported", err.Error())
}

func TestPostgresKV_Get(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"value"}).
		AddRow(
			sql.NullString{String: "bar", Valid: true},
		)

	sqlm.ExpectQuery(`SELECT value FROM kv WHERE key = \$1`).WillReturnRows(rows)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	v, err := kv.Get("k")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)
}

func TestPostgresKV_Upsert_error_exec(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectExec(`INSERT INTO kv \(key, value\) VALUES \(\$1, \$2\) ` +
		`ON CONFLICT \(key\) DO UPDATE SET value = \$2`).WillReturnError(fmt.Errorf("err1"))

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	err = kv.Upsert("k", "v")
	require.Error(t, err)
	assert.Equal(t, "error sql query, err1", err.Error())
}

func TestPostgresKV_Upsert(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectExec(`INSERT INTO kv \(key, value\) VALUES \(\$1, \$2\) ` +
		`ON CONFLICT \(key\) DO UPDATE SET value = \$2`).WillReturnResult(sqlmock.NewResult(0, 1))

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	err = kv.Upsert("k", "v")
	require.NoError(t, err)
}

func TestPostgresKV_Delete_error_exec(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectQuery(`DELETE FROM kv WHERE key = \$1`).WillReturnError(fmt.Errorf("err1"))

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	err = kv.Delete("k")
	require.Error(t, err)
	assert.Equal(t, "error sql query, err1", err.Error())
}

func TestPostgresKV_Delete(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	sqlm.ExpectQuery(`DELETE FROM kv WHERE key = \$1`).WillReturnRows(nil)

	kv := &PostgresKV{
		db:       dbx,
		tableCfg: tables.TableKV{Table: "kv", Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	err = kv.Delete("k")
	require.NoError(t, err)
}

func TestPostgresKV_CreateTable_postgres(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tableName := "kv_" + strconv.Itoa(rand.Intn(1e6))

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&sslrootcert=%s",
		"postgres",
		"secret",
		"127.0.0.1",
		35432,
		"db",
		"disable",
		"",
	)

	conn, err := sqlx.Connect("postgres", connectionString)
	require.NoError(t, err)
	defer conn.Close()

	p := &PostgresKV{
		db: conn,
		tableCfg: tables.TableKV{
			Table: tableName,
			Fields: tables.KVFields{
				Key:   "key",
				Value: "value",
			},
		},
		timeout: 0,
		logger:  zap.NewNop(),
	}

	err = p.CreateTable()
	require.NoError(t, err)

	_, err = conn.Query("SELECT key, value FROM " + tableName)
	require.NoError(t, err)
}

func TestPostgresKV_CreateTable_sqlite(t *testing.T) {
	f, err := os.CreateTemp("", "kv-")
	require.NoError(t, err)

	t.Logf("create sqlite3 file %s", f.Name())

	conn, err := sqlx.Connect("sqlite3", f.Name())
	require.NoError(t, err)
	defer conn.Close()

	tableName := "kv"

	p := &PostgresKV{
		db: conn,
		tableCfg: tables.TableKV{
			Table: tableName,
			Fields: tables.KVFields{
				Key:   "key",
				Value: "value",
			},
		},
		timeout: 0,
		logger:  zap.NewNop(),
	}

	err = p.CreateTable()
	require.NoError(t, err)

	_, err = conn.Query("SELECT key, value FROM " + tableName)
	require.NoError(t, err)
}
