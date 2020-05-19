package broccoli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	p := New()

	ss, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, ss)
}
