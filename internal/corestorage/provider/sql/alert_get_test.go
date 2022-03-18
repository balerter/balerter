package sql

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/alert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresAlert_Get(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	p, tableName := instance(t)
	defer p.db.Close()

	alertName := "foo_" + strconv.Itoa(rand.Intn(1e9))

	_, errExec := p.db.Exec("INSERT INTO " + tableName + " (name, level, count) VALUES ('" + alertName + "', 1, 2)")
	require.NoError(t, errExec)

	a, errGet := p.Get(alertName)
	require.NoError(t, errGet)
	assert.IsType(t, &alert.Alert{}, a)
	assert.Equal(t, alertName, a.Name)
	assert.Equal(t, alert.LevelSuccess, a.Level)
	assert.Equal(t, 2, a.Count)
}

func TestPostgresAlert_Get_no_rows(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	p, _ := instance(t)
	defer p.db.Close()

	alertName := "foo_" + strconv.Itoa(rand.Intn(1e9))

	a, errGet := p.Get(alertName)
	require.NoError(t, errGet)
	require.Nil(t, a)
}

func TestPostgresAlert_Get_wrong_level(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	p, tableName := instance(t)
	defer p.db.Close()

	alertName := "foo_" + strconv.Itoa(rand.Intn(1e9))

	_, errExec := p.db.Exec("INSERT INTO " + tableName + " (name, level, count) VALUES ('" + alertName + "', 99999, 2)")
	require.NoError(t, errExec)

	_, errGet := p.Get(alertName)
	require.Error(t, errGet)
	assert.Equal(t, "error parse level 99999 for alert "+alertName+", bad level", errGet.Error())
}
