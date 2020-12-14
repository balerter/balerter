package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCall_error_get_alert(t *testing.T) {
	e1 := fmt.Errorf("e1")

	em := corestorage.NewMock("m")
	em.AlertMock().On("GetOrNew", mock.Anything).Return(nil, e1)

	m := &Manager{
		storage: em,
		logger:  zap.NewNop(),
	}

	err := m.Call("", 0, "", nil)

	require.Error(t, err)
	assert.Equal(t, "error get alert , e1", err.Error())
}

func TestCall_same_level__quiet(t *testing.T) {
	mockSend := func(level, alertName, text string, channels, fields []string, image string) {
		t.Fatal("unexpected call")
	}

	a := &alert.Alert{}
	a.UpdateLevel(alert.LevelError)

	em := corestorage.NewMock("m")
	em.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)

	m := &Manager{
		storage:         em,
		logger:          zap.NewNop(),
		sendMessageFunc: mockSend,
	}

	err := m.Call("foo", alert.LevelError, "bar", &alert.Options{Quiet: true})

	require.NoError(t, err)
	assert.Equal(t, 1, a.Count())
}

func TestCall_same_level__norepeat(t *testing.T) {
	mockSend := func(level, alertName, text string, channels, fields []string, image string) {
		t.Fatal("unexpected call")
	}

	a := &alert.Alert{}
	a.UpdateLevel(alert.LevelError)

	em := corestorage.NewMock("m")
	em.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)

	m := &Manager{
		storage:         em,
		logger:          zap.NewNop(),
		sendMessageFunc: mockSend,
	}

	err := m.Call("foo", alert.LevelError, "bar", &alert.Options{Repeat: 0})

	require.NoError(t, err)
	assert.Equal(t, 1, a.Count())
}

func TestCall_same_level(t *testing.T) {
	var called bool

	mockSend := func(level, alertName, text string, channels, fields []string, image string) {
		called = true
	}

	a := &alert.Alert{}
	a.UpdateLevel(alert.LevelError)
	a.Inc()

	em := corestorage.NewMock("m")
	em.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)

	m := &Manager{
		storage:         em,
		logger:          zap.NewNop(),
		sendMessageFunc: mockSend,
	}

	err := m.Call("foo", alert.LevelError, "bar", &alert.Options{Repeat: 2})

	require.NoError(t, err)
	assert.Equal(t, 2, a.Count())

	if !called {
		t.Fatal("send has been not called")
	}
}

func TestCall_update_level__release_error(t *testing.T) {
	mockSend := func(level, alertName, text string, channels, fields []string, image string) {}

	a := &alert.Alert{}
	a.UpdateLevel(alert.LevelSuccess)
	a.Inc()

	e1 := fmt.Errorf("e1")

	em := corestorage.NewMock("m")
	em.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)
	em.AlertMock().On("Release", mock.Anything).Return(e1)

	m := &Manager{
		storage:         em,
		logger:          zap.NewNop(),
		sendMessageFunc: mockSend,
	}

	err := m.Call("foo", alert.LevelError, "bar", &alert.Options{})

	require.Error(t, err)
	assert.Equal(t, alert.LevelError, a.Level())
	assert.Equal(t, "error release alert, e1", err.Error())
}

func TestCall_update_level__quiet(t *testing.T) {
	var called bool

	mockSend := func(level, alertName, text string, channels, fields []string, image string) {
		called = true
	}

	a := &alert.Alert{}
	a.UpdateLevel(alert.LevelSuccess)
	a.Inc()

	em := corestorage.NewMock("m")
	em.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)
	em.AlertMock().On("Release", mock.Anything).Return(nil)

	m := &Manager{
		storage:         em,
		logger:          zap.NewNop(),
		sendMessageFunc: mockSend,
	}

	err := m.Call("foo", alert.LevelError, "bar", &alert.Options{Quiet: true})

	require.NoError(t, err)
	assert.Equal(t, 0, a.Count())
	assert.Equal(t, alert.LevelError, a.Level())

	if called {
		t.Fatal("unexpected call")
	}
}

func TestCall_update_level_ok(t *testing.T) {
	var called bool

	mockSend := func(level, alertName, text string, channels, fields []string, image string) {
		called = true
	}

	a := &alert.Alert{}
	a.UpdateLevel(alert.LevelSuccess)
	a.Inc()

	em := corestorage.NewMock("m")
	em.AlertMock().On("GetOrNew", mock.Anything).Return(a, nil)
	em.AlertMock().On("Release", mock.Anything).Return(nil)

	m := &Manager{
		storage:         em,
		logger:          zap.NewNop(),
		sendMessageFunc: mockSend,
	}

	err := m.Call("foo", alert.LevelError, "bar", &alert.Options{})

	require.NoError(t, err)
	assert.Equal(t, 0, a.Count())
	assert.Equal(t, alert.LevelError, a.Level())

	if !called {
		t.Fatal("unexpected not call")
	}
}
