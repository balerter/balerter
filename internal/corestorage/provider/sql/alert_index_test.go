package sql

import (
	"testing"

	"github.com/balerter/balerter/internal/alert"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresAlert_Index(t *testing.T) {
	p, tableName := instance(t)

	_ = tableName

	p.db.Exec("INSERT INTO " + tableName + " (name, level, count) VALUES ('a1', 1, 1),('a2', 2, 2),('a3', 1, 3),('a4', 1, 4)")

	alerts, errIndex := p.Index([]alert.Level{alert.LevelSuccess})
	require.NoError(t, errIndex)

	assert.Equal(t, 3, len(alerts))

	for _, a := range alerts {
		var count int
		switch a.Name {
		case "a1":
			count = 1
		case "a3":
			count = 3
		case "a4":
			count = 4
		default:
			require.True(t, false)
		}
		assert.Equal(t, count, a.Count)
	}
}
