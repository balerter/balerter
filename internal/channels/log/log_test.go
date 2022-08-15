package log

import (
	"github.com/balerter/balerter/internal/message"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"

	logConfig "github.com/balerter/balerter/internal/config/channels/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	l, err := New(logConfig.Log{Name: "test", Ignore: true}, zap.NewNop())
	require.NoError(t, err)
	assert.IsType(t, &Log{}, l)
	assert.Equal(t, "test", l.name)
	assert.True(t, l.ignore)
}

func TestLog_Name(t *testing.T) {
	l := &Log{name: "test"}
	assert.Equal(t, "test", l.Name())
}

func TestLog_Ignore(t *testing.T) {
	l := &Log{ignore: true}
	assert.True(t, l.Ignore())
}

func TestLog_Send(t *testing.T) {
	core, recordedLogs := observer.New(zapcore.InfoLevel)

	l := &Log{name: "test", logger: zap.New(core)}
	err := l.Send(&message.Message{
		Level:     "error",
		AlertName: "name1",
		Text:      "text1",
		Image:     "image1",
		Fields:    map[string]string{"foo": "bar"},
	})
	require.NoError(t, err)

	gotLog := recordedLogs.All()[0]
	assert.Equal(t, "Log channel message", gotLog.Message)
	assert.Equal(t, zap.InfoLevel, gotLog.Level)
	fields := gotLog.Context

	require.Equal(t, 6, len(fields))

	assert.Equal(t, "channel name", fields[0].Key)
	assert.Equal(t, "test", fields[0].String)

	assert.Equal(t, "alert id", fields[1].Key)
	assert.Equal(t, "name1", fields[1].String)

	assert.Equal(t, "level", fields[2].Key)
	assert.Equal(t, "error", fields[2].String)

	assert.Equal(t, "message", fields[3].Key)
	assert.Equal(t, "text1", fields[3].String)

	assert.Equal(t, "image", fields[4].Key)
	assert.Equal(t, "image1", fields[4].String)

	assert.Equal(t, "fields", fields[5].Key)
	assert.Equal(t, map[string]string{"foo": "bar"}, fields[5].Interface)
}
