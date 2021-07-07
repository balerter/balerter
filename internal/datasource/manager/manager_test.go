package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestManager_Stop(t *testing.T) {
	m1 := &modules.ModuleMock{
		StopFunc: func() error {
			return nil
		},
		NameFunc: func() string {
			return "m1"
		},
	}

	m2 := &modules.ModuleMock{
		StopFunc: func() error {
			return nil
		},
		NameFunc: func() string {
			return "m2"
		},
	}

	m := New(zap.NewNop())
	m.modules["m1"] = m1
	m.modules["m2"] = m2

	m.Stop()

	assert.Equal(t, 1, len(m1.StopCalls()))
	assert.Equal(t, 1, len(m1.NameCalls()))
	assert.Equal(t, 1, len(m2.StopCalls()))
	assert.Equal(t, 1, len(m2.NameCalls()))
}

func TestManager_Stop_Error(t *testing.T) {
	e := fmt.Errorf("error1")

	m1 := &modules.ModuleMock{
		StopFunc: func() error {
			return e
		},
		NameFunc: func() string {
			return "m1"
		},
	}

	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	m := New(logger)
	m.modules["m1"] = m1

	m.Stop()

	assert.Equal(t, 1, logs.FilterMessage("error stop module").FilterField(zap.String("name", "m1")).FilterField(zap.Error(e)).Len())

	assert.Equal(t, 1, len(m1.StopCalls()))
	assert.Equal(t, 2, len(m1.NameCalls()))
}

func TestGet(t *testing.T) {
	m := &Manager{
		modules: map[string]modules.Module{
			"m1": &modules.ModuleMock{},
			"m2": &modules.ModuleMock{},
		},
	}

	mods := m.Get()

	assert.Equal(t, 2, len(mods))
}
