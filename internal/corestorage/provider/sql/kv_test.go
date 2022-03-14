package sql

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/config/storages/core/tables"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func TestPostgresKV_All_error_query(t *testing.T) {
	tableKV := fmt.Sprintf("kv_%d", rand.Int())

	kv := &PostgresKV{
		db:       getDB(t),
		tableCfg: tables.TableKV{Table: tableKV, Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	_, err := kv.All()
	require.Error(t, err)
	assert.Equal(t, "error sql query, pq: relation \""+tableKV+"\" does not exist", err.Error())
}

func TestPostgresKV_All(t *testing.T) {
	db := getDB(t)
	tableName := fmt.Sprintf("kv_%d", rand.Int())

	kv := &PostgresKV{
		db:       db,
		tableCfg: tables.TableKV{Table: tableName, Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	errCreate := kv.CreateTable()
	require.NoError(t, errCreate)

	_, errQuery := db.Query(fmt.Sprintf(`INSERT INTO %s (key, value) VALUES ('a1', 'a2'), ('b1', 'b2')`, tableName))
	require.NoError(t, errQuery)

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

func TestPostgresKV_Put(t *testing.T) {
	db := getDB(t)
	tableName := fmt.Sprintf("kv_%d", rand.Int())

	kv := &PostgresKV{
		db:       db,
		tableCfg: tables.TableKV{Table: tableName, Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	errCreate := kv.CreateTable()
	require.NoError(t, errCreate)

	err := kv.Put("k", "v")
	require.NoError(t, err)
}

func TestPostgresKV_Get_error_no_rows(t *testing.T) {
	db := getDB(t)
	tableName := fmt.Sprintf("kv_%d", rand.Int())

	kv := &PostgresKV{
		db:       db,
		tableCfg: tables.TableKV{Table: tableName, Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	require.NoError(t, kv.CreateTable())

	_, err := kv.Get("k")
	require.Error(t, err)
	assert.Equal(t, ErrNoRow.Error(), err.Error())
}

func TestPostgresKV_Get(t *testing.T) {
	db := getDB(t)
	tableName := fmt.Sprintf("kv_%d", rand.Int())

	kv := &PostgresKV{
		db:       db,
		tableCfg: tables.TableKV{Table: tableName, Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	require.NoError(t, kv.CreateTable())

	_, errExec := db.Exec("INSERT INTO " + tableName + " VALUES ('k', 'bar')")
	require.NoError(t, errExec)

	v, err := kv.Get("k")
	require.NoError(t, err)
	assert.Equal(t, "bar", v)
}

func TestPostgresKV_Upsert(t *testing.T) {
	db := getDB(t)
	tableName := fmt.Sprintf("kv_%d", rand.Int())

	kv := &PostgresKV{
		db:       db,
		tableCfg: tables.TableKV{Table: tableName, Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	require.NoError(t, kv.CreateTable())

	_, errExec := db.Exec("INSERT INTO " + tableName + " VALUES ('k', 'bar')")
	require.NoError(t, errExec)

	err := kv.Upsert("k", "v")
	require.NoError(t, err)

	row := db.QueryRow("SELECT value FROM " + tableName + " WHERE key = 'k'")
	var v string
	errScan := row.Scan(&v)
	require.NoError(t, errScan)
	assert.Equal(t, "v", v)
}

func TestPostgresKV_Delete(t *testing.T) {
	db := getDB(t)
	tableName := fmt.Sprintf("kv_%d", rand.Int())

	kv := &PostgresKV{
		db:       db,
		tableCfg: tables.TableKV{Table: tableName, Fields: tables.KVFields{Key: "key", Value: "value"}},
		logger:   zap.NewNop(),
	}

	require.NoError(t, kv.CreateTable())

	_, errExec := db.Exec("INSERT INTO " + tableName + " VALUES ('k', 'bar')")
	require.NoError(t, errExec)

	err := kv.Delete("k")
	require.NoError(t, err)

	row := db.QueryRow("SELECT value FROM " + tableName + " WHERE key = 'k'")
	var v string
	errScan := row.Scan(&v)
	require.Error(t, errScan)
	assert.Equal(t, "sql: no rows in result set", errScan.Error())
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
