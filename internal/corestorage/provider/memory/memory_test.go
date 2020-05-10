package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory_Name(t *testing.T) {
	m := New()

	assert.Equal(t, "memory", m.Name())
}
