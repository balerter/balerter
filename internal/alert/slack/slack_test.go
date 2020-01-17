package slack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlack_Name(t *testing.T) {
	s := &Slack{name: "name1"}
	assert.Equal(t, "name1", s.Name())
}
