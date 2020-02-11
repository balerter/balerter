package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func Test_contains(t *testing.T) {
	type args struct {
		v   string
		arr []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty array",
			args: args{
				v:   "foo",
				arr: nil,
			},
			want: false,
		},
		{
			name: "not exists",
			args: args{
				v:   "foo",
				arr: []string{"bar"},
			},
			want: false,
		},
		{
			name: "exists",
			args: args{
				v:   "foo",
				arr: []string{"bar", "foo", "baz"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.v, tt.args.arr); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Send_no_channels(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := &Manager{
		logger:   logger,
		channels: map[string]alertChannel{},
	}

	m.Send(alert.LevelSuccess, "alertName", "alertText", nil, nil)

	assert.Equal(t, 0, logs.Len())
}

func TestManager_Send_skip(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	chan1 := &alertChannelMock{}

	m := &Manager{
		logger: logger,
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
	}

	m.Send(alert.LevelSuccess, "alertName", "alertText", []string{"chan2"}, nil)

	chan1.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)

	assert.Equal(t, 1, logs.FilterMessage("skip send message to channel").Len())

	chan1.AssertExpectations(t)
}

func TestManager_Send_ok(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	chan1 := &alertChannelMock{}
	chan1.On("Send", mock.Anything, mock.Anything).Return(nil)

	m := &Manager{
		logger: logger,
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
	}

	m.Send(alert.LevelSuccess, "alertName", "alertText", nil, nil)

	chan1.AssertCalled(t, "Send", mock.Anything, mock.Anything)

	assert.Equal(t, 0, logs.Len())

	chan1.AssertExpectations(t)
}

func TestManager_Send_error(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	chan1 := &alertChannelMock{}
	chan1.On("Send", mock.Anything, mock.Anything).Return(fmt.Errorf("err1"))

	m := &Manager{
		logger: logger,
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
	}

	m.Send(alert.LevelSuccess, "alertName", "alertText", nil, nil)

	chan1.AssertCalled(t, "Send", mock.Anything, mock.Anything)

	assert.Equal(t, 1, logs.FilterMessage("error send message to channel").FilterField(zap.String("channel name", "chan1")).FilterField(zap.Error(fmt.Errorf("err1"))).Len())

	chan1.AssertExpectations(t)
}
