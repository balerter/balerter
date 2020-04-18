package alert

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAcquireAlert(t *testing.T) {
	a1 := &Alert{}
	alertsPool.Put(a1)

	now := time.Now()

	a2 := AcquireAlert()

	assert.Equal(t, fmt.Sprintf("%p", a1), fmt.Sprintf("%p", a2))
	assert.Greater(t, a2.lastChange.UnixNano(), now.UnixNano())
	assert.Greater(t, a2.GetStartTime().UnixNano(), now.UnixNano())
	assert.Equal(t, LevelSuccess, a2.level)
	assert.Equal(t, 0, a2.count)
}

func TestAcquireAlert2(t *testing.T) {
	now := time.Now()

	a2 := AcquireAlert()

	assert.GreaterOrEqual(t, a2.lastChange.UnixNano(), now.UnixNano())
	assert.GreaterOrEqual(t, a2.GetStartTime().UnixNano(), now.UnixNano())
	assert.Equal(t, LevelSuccess, a2.level)
	assert.Equal(t, 0, a2.count)
}

func TestLevelFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Level
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				s: "success",
			},
			want:    LevelSuccess,
			wantErr: false,
		},
		{
			name: "warning",
			args: args{
				s: "warning",
			},
			want:    LevelWarn,
			wantErr: false,
		},
		{
			name: "warn",
			args: args{
				s: "warn",
			},
			want:    LevelWarn,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				s: "error",
			},
			want:    LevelError,
			wantErr: false,
		},
		{
			name: "wrong",
			args: args{
				s: "wrong",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LevelFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("LevelFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LevelFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	assert.Equal(t, "success", LevelSuccess.String())
	assert.Equal(t, "warning", LevelWarn.String())
	assert.Equal(t, "error", LevelError.String())

	l := Level(-1)

	assert.Panics(t, func() {
		s := l.String()
		assert.Equal(t, "", s)
	})
}
