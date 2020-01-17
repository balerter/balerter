package manager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
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

func (m *alertChannelMock) Send(alertName string, alertText string) error {
	args := m.Called(alertName, alertText)
	return args.Error(0)
}

func (m *alertChannelMock) SendSuccess(alertName string, alertText string) error {
	args := m.Called(alertName, alertText)
	return args.Error(0)
}

func (m *alertChannelMock) SendError(alertName string, alertText string) error {
	args := m.Called(alertName, alertText)
	return args.Error(0)
}

func TestManager_on(t *testing.T) {
	ch1 := &alertChannelMock{}
	ch1.On("SendError", "alert-name", "alert-text").Return(nil)

	ch2 := &alertChannelMock{}
	ch2.On("SendError", "alert-name", "alert-text").Return(nil)

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1
	m.channels["chan2"] = ch2

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	c := m.on(L)

	assert.Equal(t, 0, c)
	assert.Equal(t, 1, m.active["alert-name"])

	ch1.AssertCalled(t, "SendError", "alert-name", "alert-text")
	ch2.AssertCalled(t, "SendError", "alert-name", "alert-text")
	ch1.AssertExpectations(t)
}

func TestManager_on_without_name(t *testing.T) {
	ch1 := &alertChannelMock{}

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1

	L := lua.NewState()

	c := m.on(L)

	assert.Equal(t, 0, c)

	ch1.AssertNotCalled(t, "SendError", mock.Anything, mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_on_error(t *testing.T) {
	ch1 := &alertChannelMock{}
	ch1.On("SendError", "alert-name", "alert-text").Return(fmt.Errorf("error1"))

	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := New(logger)
	m.channels["chan1"] = ch1

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	c := m.on(L)

	require.Equal(t, 2, logs.Len())
	assert.Equal(t, 1, logs.FilterMessage("error send message to channel").FilterField(zap.String("name", "chan1")).Len())
	assert.Equal(t, 1, logs.FilterMessage("call alert ON").FilterField(zap.String("alertName", "alert-name")).Len())

	assert.Equal(t, 0, c)
	assert.Equal(t, 1, m.active["alert-name"])

	ch1.AssertCalled(t, "SendError", "alert-name", "alert-text")
	ch1.AssertExpectations(t)
}
