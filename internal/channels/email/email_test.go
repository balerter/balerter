package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailName(t *testing.T) {
	s := &Email{name: "foo"}
	assert.Equal(t, "foo", s.Name())
}
