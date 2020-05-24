package syslog

import (
	"bytes"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSend(t *testing.T) {
	s := &Syslog{}

	buf := bytes.NewBuffer([]byte{})

	s.w = buf

	mes := &message.Message{
		Level:     "foo",
		AlertName: "bar",
		Text:      "baz",
		Fields:    []string{"f1", "f2"},
		Image:     "img",
	}

	err := s.Send(mes)
	require.NoError(t, err)

	assert.Equal(t, `{"level":"foo","alert_name":"bar","text":"baz","fields":["f1","f2"],"image":"img"}`, buf.String())
}
