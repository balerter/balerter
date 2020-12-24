package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

type alertChannelMock struct {
	mock.Mock
}

func (m *alertChannelMock) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *alertChannelMock) Send(mes *message.Message) error {
	args := m.Called(mes)
	return args.Error(0)
}

func TestManager_Send_no_channels(t *testing.T) {
	m := &Manager{
		channels: map[string]alertChannel{},
	}

	err := m.Send("", "alertName", "alertText", &alert.Options{}, nil)
	require.Error(t, err)
	assert.Equal(t, ErrEmptyChannels, err)
}

func TestManager_Send_channel_not_found(t *testing.T) {
	chan1 := &alertChannelMock{}

	m := &Manager{
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
	}

	err := m.Send("", "alertName", "alertText", &alert.Options{Channels: []string{"chan2"}}, nil)
	require.Error(t, err)
	assert.Equal(t, ErrEmptyChannels, err)

	chan1.AssertNotCalled(t, "Send", mock.Anything)

	chan1.AssertExpectations(t)
}

func TestManager_Send_ok(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	chan1 := &alertChannelMock{}
	chan1.On("Send", mock.Anything).Return(nil)

	m := &Manager{
		logger: logger,
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
	}

	err := m.Send("", "alertName", "alertText", &alert.Options{}, nil)
	require.NoError(t, err)

	chan1.AssertCalled(t, "Send", mock.Anything)

	assert.Equal(t, 0, logs.Len())

	chan1.AssertExpectations(t)
}

func TestManager_Send_error(t *testing.T) {
	chan1 := &alertChannelMock{}
	chan1.On("Send", mock.Anything).Return(fmt.Errorf("err1"))

	m := &Manager{
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
	}

	err := m.Send("", "alertName", "alertText", &alert.Options{}, nil)
	require.NoError(t, err)

	chan1.AssertCalled(t, "Send", mock.Anything)

	chan1.AssertExpectations(t)
}
