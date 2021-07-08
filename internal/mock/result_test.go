package mock

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleMock_Result(t *testing.T) {
	r := &registryMock{}
	r.On("Result").Return([]modules.TestResult{{
		ScriptName: "a",
		ModuleName: "b",
		Message:    "c",
		Ok:         true,
	}})

	mm := &ModuleMock{
		name:     "m",
		errors:   []string{"foo"},
		registry: r,
	}

	res, err := mm.Result()
	require.NoError(t, err)

	require.Equal(t, 2, len(res))
	assert.Equal(t, "foo", res[0].Message)
	assert.Equal(t, "m", res[0].ModuleName)
	assert.Equal(t, "", res[0].ScriptName)
	assert.False(t, res[0].Ok)

	assert.Equal(t, "c", res[1].Message)
	assert.Equal(t, "m", res[1].ModuleName)
	assert.Equal(t, "a", res[1].ScriptName)
	assert.True(t, res[1].Ok)
}
