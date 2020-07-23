package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	rnr := New(nil, nil, nil, nil, nil)
	assert.IsType(t, &Runner{}, rnr)
}
