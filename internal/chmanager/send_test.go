package manager

import (
	"github.com/balerter/balerter/internal/message"
	"github.com/stretchr/testify/mock"
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

//func TestManager_Send_no_channels(t *testing.T) {
//	m := &ChannelsManager{
//		channels: map[string]alertChannel{},
//	}
//
//	err := m.Send(alert.New("alertName"), "alertText", &alert.Options{})
//	require.Error(t, err)
//	assert.Equal(t, ErrEmptyChannels, err)
//}
//
//func TestManager_Send_channel_not_found(t *testing.T) {
//	chan1 := &alertChannelMock{}
//
//	m := &ChannelsManager{
//		channels: map[string]alertChannel{
//			"chan1": chan1,
//		},
//	}
//
//	err := m.Send("", "alertName", "alertText", &alert.Options{Channels: []string{"chan2"}}, nil)
//	require.Error(t, err)
//	assert.Equal(t, ErrEmptyChannels, err)
//
//	chan1.AssertNotCalled(t, "Send", mock.Anything)
//
//	chan1.AssertExpectations(t)
//}
//
//func TestManager_Send_ok(t *testing.T) {
//	core, logs := observer.New(zap.DebugLevel)
//	logger := zap.New(core)
//
//	chan1 := &alertChannelMock{}
//	chan1.On("Send", mock.Anything).Return(nil)
//
//	m := &ChannelsManager{
//		logger: logger,
//		channels: map[string]alertChannel{
//			"chan1": chan1,
//		},
//	}
//
//	err := m.Send("", "alertName", "alertText", &alert.Options{}, nil)
//	require.NoError(t, err)
//
//	chan1.AssertCalled(t, "Send", mock.Anything)
//
//	assert.Equal(t, 0, logs.Len())
//
//	chan1.AssertExpectations(t)
//}
//
//func TestManager_Send_error(t *testing.T) {
//	chan1 := &alertChannelMock{}
//	chan1.On("Send", mock.Anything).Return(fmt.Errorf("err1"))
//
//	m := &ChannelsManager{
//		channels: map[string]alertChannel{
//			"chan1": chan1,
//		},
//	}
//
//	err := m.Send("", "alertName", "alertText", &alert.Options{}, nil)
//	require.NoError(t, err)
//
//	chan1.AssertCalled(t, "Send", mock.Anything)
//
//	chan1.AssertExpectations(t)
//}
