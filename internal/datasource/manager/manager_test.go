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
	m1 := &modules.ModuleMock{}
	m1.On("Stop").Return(nil)
	m1.On("Name").Return("m1")

	m2 := &modules.ModuleMock{}
	m2.On("Stop").Return(nil)
	m2.On("Name").Return("m2")

	m := New(zap.NewNop())
	m.modules["m1"] = m1
	m.modules["m2"] = m2

	m.Stop()

	m1.AssertCalled(t, "Stop")
	m2.AssertCalled(t, "Stop")

	m1.AssertExpectations(t)
	m2.AssertExpectations(t)
}

func TestManager_Stop_Error(t *testing.T) {
	e := fmt.Errorf("error1")

	m1 := &modules.ModuleMock{}
	m1.On("Stop").Return(e)
	m1.On("Name").Return("m1")

	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	m := New(logger)
	m.modules["m1"] = m1

	m.Stop()

	assert.Equal(t, 1, logs.FilterMessage("error stop module").FilterField(zap.String("name", "m1")).FilterField(zap.Error(e)).Len())

	m1.AssertCalled(t, "Stop")

	m1.AssertExpectations(t)
}
