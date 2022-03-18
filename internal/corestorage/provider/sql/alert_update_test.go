package sql

import (
	"testing"

	"github.com/balerter/balerter/internal/alert"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresAlert_Update(t *testing.T) {
	p, tableName := instance(t)

	p.db.Exec("INSERT INTO " + tableName + " (name, level, count) VALUES ('foo', 1, 10)")

	a, ok, errUpdate := p.Update("foo", alert.LevelError)
	require.NoError(t, errUpdate)
	assert.True(t, ok)

	assert.Equal(t, a.Level, alert.LevelError)

	row := p.db.QueryRow("SELECT level, count FROM " + tableName + " WHERE name = 'foo'")
	var level int
	var count int
	errScan := row.Scan(&level, &count)
	require.NoError(t, errScan)
	assert.Equal(t, int(alert.LevelError), level)
	assert.Equal(t, 1, count)
}

func TestPostgresAlert_Update_new_alert(t *testing.T) {
	p, tableName := instance(t)

	a, ok, errUpdate := p.Update("foo", alert.LevelError)
	require.NoError(t, errUpdate)
	assert.True(t, ok)

	assert.Equal(t, a.Level, alert.LevelError)

	row := p.db.QueryRow("SELECT level, count FROM " + tableName + " WHERE name = 'foo'")
	var level int
	var count int
	errScan := row.Scan(&level, &count)
	require.NoError(t, errScan)
	assert.Equal(t, int(alert.LevelError), level)
	assert.Equal(t, 1, count)
}
