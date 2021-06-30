package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory_Name(t *testing.T) {
	m := New()

	assert.Equal(t, "memory", m.Name())
}

func TestMemory_KV(t *testing.T) {
	kv := &storageKV{}
	m := Memory{
		kv: kv,
	}

	assert.Equal(t, kv, m.KV())
}

func TestMemory_Alert(t *testing.T) {
	a := &storageAlert{}
	m := Memory{
		alert: a,
	}

	assert.Equal(t, a, m.Alert())
}

func TestMemory_Stop(t *testing.T) {
	m := Memory{}
	assert.NoError(t, m.Stop())
}
