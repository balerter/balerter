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

func TestManager_off(t *testing.T) {
	ch1 := &alertChannelMock{}
	ch1.On("SendSuccess", "alert-name", "alert-text").Return(nil)

	ch2 := &alertChannelMock{}
	ch2.On("SendSuccess", "alert-name", "alert-text").Return(nil)

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1
	m.channels["chan2"] = ch2
	m.active["alert-name"] = &alertInfo{count: 10}

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	c := m.off(L)

	assert.Equal(t, 0, c)
	_, ok := m.active["alert-name"]
	assert.False(t, ok)

	ch1.AssertCalled(t, "SendSuccess", "alert-name", "alert-text")
	ch2.AssertCalled(t, "SendSuccess", "alert-name", "alert-text")
	ch1.AssertExpectations(t)
}

func TestManager_off_without_name(t *testing.T) {
	ch1 := &alertChannelMock{}

	m := New(zap.NewNop())
	m.channels["chan1"] = ch1

	L := lua.NewState()

	c := m.off(L)

	assert.Equal(t, 0, c)

	ch1.AssertNotCalled(t, "SendSuccess", mock.Anything, mock.Anything)
	ch1.AssertExpectations(t)
}

func TestManager_off_error(t *testing.T) {
	ch1 := &alertChannelMock{}
	ch1.On("SendSuccess", "alert-name", "alert-text").Return(fmt.Errorf("error1"))

	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	m := New(logger)
	m.channels["chan1"] = ch1
	m.active["alert-name"] = &alertInfo{count: 10}

	L := lua.NewState()
	L.Push(lua.LString("alert-name"))
	L.Push(lua.LString("alert-text"))

	c := m.off(L)

	require.Equal(t, 2, logs.Len())
	assert.Equal(t, 1, logs.FilterMessage("error send message to channel").FilterField(zap.String("name", "chan1")).Len())
	assert.Equal(t, 1, logs.FilterMessage("call alert OFF").FilterField(zap.String("alertName", "alert-name")).Len())

	assert.Equal(t, 0, c)
	_, ok := m.active["alert-name"]
	assert.False(t, ok)

	ch1.AssertCalled(t, "SendSuccess", "alert-name", "alert-text")
	ch1.AssertExpectations(t)

	//assert.True(t, bytes.Contains(sink.Bytes(), []byte(`"msg":"error send message to channel","name":"chan1","error":"error1"`)))
}
