package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := New()
	assert.Equal(t, defaultScriptsUpdateInterval, cfg.Scripts.Sources.UpdateInterval)
}
