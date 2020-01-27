package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/message"
	"github.com/balerter/balerter/internal/script/script"
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

func (m *alertChannelMock) Send(message *message.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *alertChannelMock) SendSuccess(message *message.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *alertChannelMock) SendError(message *message.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestManager_on(t *testing.T) {
	ch1 := &alertChannelMock{}
	ch1.On("SendError", message.New("alert-name", "alert-text")).Return(nil)

	ch2 := &alertChannelMock{}
	ch2.On("SendError", message.New("alert-name", "alert-text")).Return(nil)

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1
	m.channels["chan2"] = ch2

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	f := m.on(&script.Script{})
	c := f(L)

	assert.Equal(t, 0, c)
	i, ok := m.active["alert-name"]
	require.True(t, ok)
	assert.True(t, i.Active)

	ch1.AssertCalled(t, "SendError", mock.Anything)
	ch2.AssertCalled(t, "SendError", mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_on_without_name(t *testing.T) {
	ch1 := &alertChannelMock{}

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1

	L := lua.NewState()

	f := m.on(&script.Script{})
	c := f(L)

	assert.Equal(t, 1, c)
	err := L.Get(1)
	assert.Equal(t, "error: alertName must be specified and not empty", err.String())

	ch1.AssertNotCalled(t, "SendError", mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_on_error(t *testing.T) {
	ch1 := &alertChannelMock{}
	ch1.On("SendError", message.New("alert-name", "alert-text")).Return(fmt.Errorf("error1"))

	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := New(logger)
	m.channels["chan1"] = ch1

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	f := m.on(&script.Script{})
	c := f(L)

	require.Equal(t, 2, logs.Len())
	assert.Equal(t, 1, logs.FilterMessage("error send message to channel").FilterField(zap.String("name", "chan1")).Len())
	assert.Equal(t, 1, logs.FilterMessage("call alert ON").FilterField(zap.String("alertName", "alert-name")).Len())

	assert.Equal(t, 0, c)
	i, ok := m.active["alert-name"]
	require.True(t, ok)
	assert.True(t, i.Active)

	ch1.AssertCalled(t, "SendError", mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_off(t *testing.T) {
	ch1 := &alertChannelMock{}

	ch2 := &alertChannelMock{}

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1
	m.channels["chan2"] = ch2
	m.active["alert-name"] = &alertInfo{}

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	f := m.off(&script.Script{})
	c := f(L)

	assert.Equal(t, 0, c)
	info, ok := m.active["alert-name"]
	assert.True(t, ok)
	assert.False(t, info.Active)

	ch1.AssertNotCalled(t, "SendSuccess", mock.Anything)
	ch2.AssertNotCalled(t, "SendSuccess", mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_off_without_name(t *testing.T) {
	ch1 := &alertChannelMock{}

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1

	L := lua.NewState()

	f := m.off(&script.Script{})
	c := f(L)

	assert.Equal(t, 1, c)
	err := L.Get(1)
	assert.Equal(t, "error: alertName must be specified and not empty", err.String())

	ch1.AssertNotCalled(t, "SendSuccess", mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_off_error(t *testing.T) {
	ch1 := &alertChannelMock{}
	//ch1.On("SendSuccess", message.New("alert-name", "alert-text")).Return(fmt.Errorf("error1"))

	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := New(logger)
	m.channels["chan1"] = ch1
	m.active["alert-name"] = &alertInfo{}

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	f := m.off(&script.Script{})
	c := f(L)

	require.Equal(t, 1, logs.Len())
	//assert.Equal(t, 1, logs.FilterMessage("error send message to channel").FilterField(zap.String("name", "chan1")).Len())
	assert.Equal(t, 1, logs.FilterMessage("call alert OFF").FilterField(zap.String("alertName", "alert-name")).Len())

	assert.Equal(t, 0, c)
	info, ok := m.active["alert-name"]
	assert.True(t, ok)
	assert.False(t, info.Active)

	ch1.AssertNotCalled(t, "SendSuccess", mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_getArgs_with_fields(t *testing.T) {
	m := &Manager{}

	L := lua.NewState()
	L.Push(lua.LString("alert1"))
	L.Push(lua.LString("text1"))

	fs := &lua.LTable{}

	fs.Append(lua.LString("f1"))

	L.Push(fs)

	alertName, text, fields, err := m.getArgs(L)
	require.NoError(t, err)
	assert.Equal(t, "alert1", alertName)
	assert.Equal(t, "text1", text)

	require.Equal(t, 1, len(fields))
	assert.Equal(t, "f1", fields[0])
}

func TestManager_getArgs_without_fields(t *testing.T) {
	m := &Manager{}

	L := lua.NewState()
	L.Push(lua.LString("alert1"))
	L.Push(lua.LString("text1"))

	alertName, text, fields, err := m.getArgs(L)
	require.NoError(t, err)
	assert.Equal(t, "alert1", alertName)
	assert.Equal(t, "text1", text)

	require.Equal(t, 0, len(fields))
}
