package manager

//func TestManager_Send(t *testing.T) {
//	ch1 := &alertChannelMock{}
//	ch2 := &alertChannelMock{}
//
//	ch1.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)
//
//	core, logs := observer.New(zap.DebugLevel)
//	logger := zap.New(core)
//
//	m := New(logger)
//	m.channels["ch1"] = ch1
//	m.channels["ch2"] = ch2
//
//	err := m.Send("message", []string{"ch1"})
//	require.NoError(t, err)
//
//	require.Equal(t, 0, logs.Len())
//
//	ch1.AssertCalled(t, "Send", mock.Anything)
//	ch2.AssertNotCalled(t, "Send", mock.Anything)
//
//	logger.Check(zap.ErrorLevel, "")
//
//	ch1.AssertExpectations(t)
//	ch2.AssertExpectations(t)
//}
//
//func TestManager_Send_WrongChannel(t *testing.T) {
//	ch1 := &alertChannelMock{}
//
//	core, logs := observer.New(zap.DebugLevel)
//	logger := zap.New(core)
//
//	m := New(logger)
//	m.channels["ch1"] = ch1
//
//	err := m.Send("message", []string{"ch2"})
//	require.NoError(t, err)
//
//	require.Equal(t, 1, logs.Len())
//	assert.Equal(t, 1, logs.FilterMessage("channel not found").FilterField(zap.String("channel", "ch2")).Len())
//
//	ch1.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
//
//	ch1.AssertExpectations(t)
//}
//
//func TestManager_Send_error(t *testing.T) {
//	e := fmt.Errorf("error")
//
//	ch1 := &alertChannelMock{}
//	ch1.On("Send", mock.Anything).Return(e)
//
//	core, logs := observer.New(zap.DebugLevel)
//	logger := zap.New(core)
//
//	m := New(logger)
//	m.channels["ch1"] = ch1
//
//	err := m.Send("message", []string{"ch1"})
//	require.NoError(t, err)
//
//	require.Equal(t, 1, logs.Len())
//	assert.Equal(t, 1, logs.FilterMessage("error send message").FilterField(zap.String("channel", "ch1")).FilterField(zap.Error(e)).Len())
//
//	ch1.AssertCalled(t, "Send", mock.Anything)
//
//	logger.Check(zap.ErrorLevel, "")
//
//	ch1.AssertExpectations(t)
//}
