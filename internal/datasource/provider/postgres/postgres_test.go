package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/balerter/balerter/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNew_ErrorConnect(t *testing.T) {
	mockConnFunc := func(string, string) (*sqlx.DB, error) {
		return nil, fmt.Errorf("err1")
	}

	cfg := &config.DataSourcePostgres{}

	_, err := New(cfg, mockConnFunc, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

func TestNew_ErrorPing(t *testing.T) {
	db, dbmock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	mockConnFunc := func(string, string) (*sqlx.DB, error) {
		return sqlx.NewDb(db, "sqlmock"), nil
	}

	dbmock.ExpectPing().WillReturnError(fmt.Errorf("err2"))

	cfg := &config.DataSourcePostgres{}

	_, err = New(cfg, mockConnFunc, zap.NewNop())

	require.Error(t, err)
	assert.Equal(t, "err2", err.Error())
}

func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)

	mockConnFunc := func(string, string) (*sqlx.DB, error) {
		return sqlx.NewDb(db, "sqlmock"), nil
	}

	cfg := &config.DataSourcePostgres{}

	p, err := New(cfg, mockConnFunc, zap.NewNop())

	require.NoError(t, err)
	assert.IsType(t, &Postgres{}, p)
}
