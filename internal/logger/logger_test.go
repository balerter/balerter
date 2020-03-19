package logger

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func Test_convertStringLevelToZap(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want zapcore.Level
	}{
		{
			name: "ERROR",
			args: args{level: "ERROR"},
			want: zap.ErrorLevel,
		},
		{
			name: "WARN",
			args: args{level: "WARN"},
			want: zap.WarnLevel,
		},
		{
			name: "INFO",
			args: args{level: "INFO"},
			want: zap.InfoLevel,
		},
		{
			name: "DEBUG",
			args: args{level: "DEBUG"},
			want: zap.DebugLevel,
		},
		{
			name: "bad",
			args: args{level: "bad"},
			want: zap.ErrorLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertStringLevelToZap(tt.args.level); got != tt.want {
				t.Errorf("convertStringLevelToZap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew_NoDebug(t *testing.T) {
	l, err := New("ERROR", false)
	require.NoError(t, err)
	assert.IsType(t, &Logger{}, l)
	assert.IsType(t, &zap.Logger{}, l.logger)
}

func TestNew_Debug(t *testing.T) {
	l, err := New("ERROR", true)
	require.NoError(t, err)
	assert.IsType(t, &Logger{}, l)
	assert.IsType(t, &zap.Logger{}, l.logger)
}
