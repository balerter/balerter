package clickhouse

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var _db *sqlx.DB

func getDB(t *testing.T) *sqlx.DB {
	if _db != nil {
		return _db
	}

	connString := fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s&%s",
		"127.0.0.1",
		9000,
		"default",
		"",
		"default",
		"",
	)

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ctxCancel()

	var err error

	_db, err = sqlx.ConnectContext(ctx, "clickhouse", connString)
	require.NoError(t, err)
	return _db
}
