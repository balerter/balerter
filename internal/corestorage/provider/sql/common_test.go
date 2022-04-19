package sql

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/config/storages/core/tables"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func instance(t *testing.T) (*PostgresAlert, string) {
	rand.Seed(time.Now().UnixNano())

	alertTableName := "alert_" + strconv.Itoa(rand.Intn(1e9))

	db, errOpen := sqlx.Open("postgres", "postgres://postgres:secret@127.0.0.1:35432/postgres?sslmode=disable")
	require.NoError(t, errOpen)

	p := &PostgresAlert{
		logger: zap.NewNop(),
		db:     db,
		tableCfg: tables.TableAlerts{
			Table: alertTableName,
			Fields: tables.AlertFields{
				Name:      "name",
				Level:     "level",
				Count:     "count",
				UpdatedAt: "updated_at",
				CreatedAt: "created_at",
			},
			CreateTable: false,
		},
	}

	errCreateTable := p.CreateTable()
	require.NoError(t, errCreateTable)

	return p, alertTableName
}
