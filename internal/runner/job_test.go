package runner

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
	"time"
)

func TestRunner_createLuaState(t *testing.T) {
	m1 := &moduleMock{
		NameFunc: func() string {
			return "module1"
		},
		GetLoaderFunc: func(_ *script.Script) lua.LGFunction {
			return func() lua.LGFunction {
				return func(state *lua.LState) int {
					return 0
				}
			}()
		},
	}

	dsManager := &dsManagerMock{
		GetFunc: func() []modules.Module {
			return []modules.Module{m1}
		},
	}

	storagesManager := &storagesManagerMock{
		GetFunc: func() []modules.Module {
			return []modules.Module{m1}
		},
	}

	alertManager := &moduleMock{
		NameFunc: func() string {
			return "alert1"
		},
		GetLoaderFunc: func(_ *script.Script) lua.LGFunction {
			return func() lua.LGFunction {
				return func(state *lua.LState) int {
					return 0
				}
			}()
		},
	}

	rnr := &Runner{
		logger:          zap.NewNop(),
		dsManager:       dsManager,
		storagesManager: storagesManager,
		coreModules:     []modules.Module{alertManager},
	}

	j := &Job{name: "job1"}

	err := rnr.createLuaState(j, nil)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(m1.NameCalls()))
	assert.Equal(t, 2, len(m1.GetLoaderCalls()))

	require.NotNil(t, j.luaState)
}

func TestNew(t *testing.T) {
	j := newJob(&script.Script{}, zap.NewNop())
	assert.IsType(t, &Job{}, j)
}

func TestJob_Stop(t *testing.T) {
	j := &Job{
		luaState: lua.NewState(),
	}

	j.Stop()
}

func TestJob_Run(t *testing.T) {
	j := &Job{
		luaState: lua.NewState(),
		script:   &script.Script{Timeout: time.Millisecond * 100, Body: []byte("print(1)")},
		logger:   zap.NewNop(),
	}

	j.Run()
}

func TestJob_Run_already_running(t *testing.T) {
	core, logger := observer.New(zap.DebugLevel)

	j := &Job{
		running: 1,
		logger:  zap.New(core),
	}

	j.Run()

	assert.Equal(t, 1, logger.FilterMessage("job already running").Len())
}
