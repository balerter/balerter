package alert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestLevelFromInt(t *testing.T) {
	type args struct {
		i int
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
				i: 1,
			},
			want:    LevelSuccess,
			wantErr: false,
		},
		{
			name: "warn",
			args: args{
				i: 2,
			},
			want:    LevelWarn,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				i: 3,
			},
			want:    LevelError,
			wantErr: false,
		},
		{
			name: "bad 1",
			args: args{
				i: 0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "bad 2",
			args: args{
				i: 4,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LevelFromInt(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("LevelFromInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LevelFromInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelNumString(t *testing.T) {
	assert.Equal(t, "1", LevelSuccess.NumString())
	assert.Equal(t, "2", LevelWarn.NumString())
	assert.Equal(t, "3", LevelError.NumString())
}

func TestNew(t *testing.T) {
	a := New("foo")
	assert.IsType(t, &Alert{}, a)
	assert.Equal(t, "foo", a.Name)
	assert.Equal(t, LevelSuccess, a.Level)
}
