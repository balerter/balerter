package syslog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	s := &Email{name: "foo"}
	assert.Equal(t, "foo", s.Name())
}
