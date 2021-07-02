package sql

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestPostgresAlert_CreateTable_postgres(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tableName := "alert_" + strconv.Itoa(rand.Intn(1e6))

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

	p := &PostgresAlert{
		db: conn,
		tableCfg: tables.TableAlerts{
			Table: tableName,
			Fields: tables.AlertFields{
				Name:      "id",
				Level:     "level",
				Count:     "count",
				UpdatedAt: "updated_at",
				CreatedAt: "created_at",
			},
		},
		timeout: 0,
		logger:  zap.NewNop(),
	}

	err = p.CreateTable()
	require.NoError(t, err)

	_, err = conn.Query("SELECT id, level, count, updated_at, created_at FROM " + tableName)
	require.NoError(t, err)
}

func TestPostgresAlert_CreateTable_sqlite(t *testing.T) {
	f, err := os.CreateTemp("", "alert-")
	require.NoError(t, err)

	t.Logf("create sqlite3 file %s", f.Name())

	conn, err := sqlx.Connect("sqlite3", f.Name())
	require.NoError(t, err)
	defer conn.Close()

	tableName := "kv"

	p := &PostgresAlert{
		db: conn,
		tableCfg: tables.TableAlerts{
			Table: tableName,
			Fields: tables.AlertFields{
				Name:      "id",
				Level:     "level",
				Count:     "count",
				UpdatedAt: "updated_at",
				CreatedAt: "created_at",
			},
		},
		timeout: 0,
		logger:  zap.NewNop(),
	}

	err = p.CreateTable()
	require.NoError(t, err)

	_, err = conn.Query("SELECT id, level, count, updated_at, created_at FROM " + tableName)
	require.NoError(t, err)
}
