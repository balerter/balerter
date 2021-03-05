package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	core, logger := observer.New(zap.DebugLevel)
	m := &ChannelsManager{
		channels: map[string]alertChannel{},
		logger:   zap.New(core),
	}

	m.Send(alert.New("alertName"), "alertText", &alert.Options{})

	assert.Equal(t, 1, logger.FilterMessage("the message was not sent, empty channels").Len())
}

func TestManager_Send_channel_not_found(t *testing.T) {
	core, logger := observer.New(zap.DebugLevel)
	chan1 := &alertChannelMock{}

	m := &ChannelsManager{
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
		logger: zap.New(core),
	}

	m.Send(alert.New("alertName"), "alertText", &alert.Options{Channels: []string{"chan2"}})

	chan1.AssertNotCalled(t, "Send", mock.Anything)

	chan1.AssertExpectations(t)

	assert.Equal(t, 1, logger.FilterMessage("the message was not sent, empty channels").Len())
	assert.Equal(t, 1, logger.FilterMessage("channel not found").Len())
}

func TestManager_Send_ok(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	chan1 := &alertChannelMock{}
	chan1.On("Send", mock.Anything).Return(nil)

	m := &ChannelsManager{
		logger: logger,
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
	}

	m.Send(alert.New("alertName"), "alertText", &alert.Options{})

	chan1.AssertCalled(t, "Send", mock.Anything)

	assert.Equal(t, 0, logs.Len())

	chan1.AssertExpectations(t)
}

func TestManager_Send_error(t *testing.T) {
	chan1 := &alertChannelMock{}
	chan1.On("Send", mock.Anything).Return(fmt.Errorf("err1"))

	core, logger := observer.New(zap.DebugLevel)
	m := &ChannelsManager{
		channels: map[string]alertChannel{
			"chan1": chan1,
		},
		logger: zap.New(core),
	}

	m.Send(alert.New("alertName"), "alertText", &alert.Options{})

	chan1.AssertCalled(t, "Send", mock.Anything)

	chan1.AssertExpectations(t)

	assert.Equal(t, 1, logger.FilterMessage("error send the message to the channel").Len())
}
