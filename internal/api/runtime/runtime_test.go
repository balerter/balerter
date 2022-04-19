package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	r := New(nil, nil)
	assert.IsType(t, &Runtime{}, r)
}
