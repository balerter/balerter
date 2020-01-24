package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := New()
	assert.IsType(t, &Config{}, cfg)
}
