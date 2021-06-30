package memory

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorageAlert_Get_not_found(t *testing.T) {
	a := &storageAlert{
		alerts: map[string]*alert.Alert{},
	}

	_, err := a.Get("k")
	require.Error(t, err)
	assert.Equal(t, "alert not found", err.Error())
}

func TestStorageAlert_Get(t *testing.T) {
	al := &alert.Alert{}
	a := &storageAlert{
		alerts: map[string]*alert.Alert{"k": al},
	}

	ae, err := a.Get("k")
	require.NoError(t, err)
	assert.Equal(t, al, ae)
}

func TestStorageAlert_Index(t *testing.T) {
	a1 := &alert.Alert{Level: alert.LevelSuccess}
	a2 := &alert.Alert{Level: alert.LevelWarn}
	a3 := &alert.Alert{Level: alert.LevelError}
	a := &storageAlert{
		alerts: map[string]*alert.Alert{"a1": a1, "a2": a2, "a3": a3},
	}

	data, err := a.Index([]alert.Level{alert.LevelSuccess, alert.LevelWarn})
	require.NoError(t, err)
	assert.Equal(t, 2, len(data))

	assert.Equal(t, a1, data[0])
	assert.Equal(t, a2, data[1])
}

func TestStorageAlert_Update_no_alert_success(t *testing.T) {
	a1 := &alert.Alert{Level: alert.LevelSuccess}
	a2 := &alert.Alert{Level: alert.LevelWarn}
	a3 := &alert.Alert{Level: alert.LevelError}
	a := &storageAlert{
		alerts: map[string]*alert.Alert{"a1": a1, "a2": a2, "a3": a3},
	}

	ae, updated, err := a.Update("a4", alert.LevelSuccess)
	require.NoError(t, err)
	assert.False(t, updated)
	assert.Equal(t, alert.LevelSuccess, ae.Level)
	assert.Equal(t, "a4", ae.Name)
	assert.Equal(t, 0, ae.Count)
}

func TestStorageAlert_Update_no_alert_not_success(t *testing.T) {
	a1 := &alert.Alert{Level: alert.LevelSuccess}
	a2 := &alert.Alert{Level: alert.LevelWarn}
	a3 := &alert.Alert{Level: alert.LevelError}
	a := &storageAlert{
		alerts: map[string]*alert.Alert{"a1": a1, "a2": a2, "a3": a3},
	}

	ae, updated, err := a.Update("a4", alert.LevelWarn)
	require.NoError(t, err)
	assert.True(t, updated)
	assert.Equal(t, alert.LevelWarn, ae.Level)
	assert.Equal(t, "a4", ae.Name)
	assert.Equal(t, 0, ae.Count)
}

func TestStorageAlert_Update_exists_not_change_level(t *testing.T) {
	a1 := &alert.Alert{Name: "a1", Level: alert.LevelSuccess}
	a2 := &alert.Alert{Name: "a2", Level: alert.LevelWarn}
	a3 := &alert.Alert{Name: "a3", Level: alert.LevelError}
	a := &storageAlert{
		alerts: map[string]*alert.Alert{"a1": a1, "a2": a2, "a3": a3},
	}

	ae, updated, err := a.Update("a2", alert.LevelWarn)
	require.NoError(t, err)
	assert.False(t, updated)
	assert.Equal(t, alert.LevelWarn, ae.Level)
	assert.Equal(t, "a2", ae.Name)
	assert.Equal(t, 1, ae.Count)
}

func TestStorageAlert_Update_exists_change_level(t *testing.T) {
	a1 := &alert.Alert{Name: "a1", Level: alert.LevelSuccess}
	a2 := &alert.Alert{Name: "a2", Level: alert.LevelWarn}
	a3 := &alert.Alert{Name: "a3", Level: alert.LevelError}
	a := &storageAlert{
		alerts: map[string]*alert.Alert{"a1": a1, "a2": a2, "a3": a3},
	}

	ae, updated, err := a.Update("a2", alert.LevelSuccess)
	require.NoError(t, err)
	assert.True(t, updated)
	assert.Equal(t, alert.LevelSuccess, ae.Level)
	assert.Equal(t, "a2", ae.Name)
	assert.Equal(t, 0, ae.Count)
}
