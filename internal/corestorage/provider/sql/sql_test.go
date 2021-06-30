package sql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNew_bad_connection(t *testing.T) {
	_, err := New("foo", "sqlmock", "sqlmock", tables.TableAlerts{}, tables.TableKV{}, time.Second, zap.NewNop())
	require.Error(t, err)
	assert.Equal(t, "expected a connection to be available, but it is not", err.Error())
}

func TestName(t *testing.T) {
	s := &SQL{name: "foo"}
	assert.Equal(t, "foo", s.Name())
}

func TestKV(t *testing.T) {
	kv := &PostgresKV{}
	s := &SQL{kv: kv}
	assert.Equal(t, kv, s.KV())
}

func TestAlert(t *testing.T) {
	a := &PostgresAlert{}
	s := &SQL{alerts: a}
	assert.Equal(t, a, s.Alert())
}

func TestStop_error(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")
	sqlm.ExpectClose().WillReturnError(fmt.Errorf("err1"))

	s := &SQL{db: dbx}
	err = s.Stop()
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestStop(t *testing.T) {
	db, sqlm, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")
	sqlm.ExpectClose()

	s := &SQL{db: dbx}
	err = s.Stop()
	require.NoError(t, err)
}
