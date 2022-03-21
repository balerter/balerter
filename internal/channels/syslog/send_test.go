package syslog

import (
	"bytes"
	"testing"

	"github.com/balerter/balerter/internal/message"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	s := &Syslog{}

	buf := bytes.NewBuffer([]byte{})

	s.w = buf

	mes := &message.Message{
		Level:     "foo",
		AlertName: "bar",
		Text:      "baz",
		Image:     "img",
	}

	err := s.Send(mes)
	require.NoError(t, err)

	assert.Equal(t, `{"level":"foo","alert_name":"bar","text":"baz","image":"img"}`, buf.String())
}
