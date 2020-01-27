package manager

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestManager_GetAlerts(t *testing.T) {
	m := &Manager{}

	m.alerts = map[string]*alertInfo{
		"a1": {Active: true, ScriptName: "s1"},
		"a2": {Active: false, ScriptName: "s2"},
		"a3": {Active: true, ScriptName: "s3"},
	}

	data := m.GetAlerts()
	require.Equal(t, 3, len(data))

	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	assert.Equal(t, "a1", data[0].Name)
	assert.Equal(t, "s1", data[0].ScriptName)
	assert.Equal(t, true, data[0].Active)

	assert.Equal(t, "a2", data[1].Name)
	assert.Equal(t, "s2", data[1].ScriptName)
	assert.Equal(t, false, data[1].Active)

	assert.Equal(t, "a3", data[2].Name)
	assert.Equal(t, "s3", data[2].ScriptName)
	assert.Equal(t, true, data[2].Active)
}
