package sql

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type PostgresAlert struct {
	db      *sqlx.DB
	table   string
	timeout time.Duration
	logger  *zap.Logger
}
