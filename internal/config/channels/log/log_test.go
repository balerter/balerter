package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLog_Validate_error(t *testing.T) {
	l := &Log{Name: ""}

	err := l.Validate()
	require.Error(t, err)
	assert.Equal(t, "name must be not empty", err.Error())
}

func TestLog_Validate(t *testing.T) {
	l := &Log{Name: "foo"}

	err := l.Validate()
	require.NoError(t, err)
}
