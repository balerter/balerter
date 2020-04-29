package runner

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRunner_Watch(t *testing.T) {
	scripts := make([]*script.Script, 0)

	scriptsMgr := &scriptManagerMock{}
	scriptsMgr.On("Get").Return(scripts, nil)

	ctx, ctxCancel := context.WithCancel(context.Background())
	var wg *sync.WaitGroup

	rnr := &Runner{
		pool:           make(map[string]*runningJob),
		scriptsManager: scriptsMgr,
		logger:         zap.NewNop(),
	}

	time.AfterFunc(time.Millisecond*200, func() {
		ctxCancel()
	})

	var cancelled bool

	go func() {
		rnr.Watch(ctx, ctxCancel, wg, false)
		cancelled = true
	}()

	<-time.After(time.Millisecond * 500)
	assert.True(t, cancelled)

	scriptsMgr.AssertCalled(t, "Get")
	scriptsMgr.AssertExpectations(t)
}

func TestRunner_Watch_Error(t *testing.T) {
	scripts := make([]*script.Script, 0)

	e := fmt.Errorf("error1")

	scriptsMgr := &scriptManagerMock{}
	scriptsMgr.On("Get").Return(scripts, e)

	ctx, ctxCancel := context.WithCancel(context.Background())
	var wg *sync.WaitGroup

	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	rnr := &Runner{
		updateInterval: time.Second,
		pool:           make(map[string]*runningJob),
		scriptsManager: scriptsMgr,
		logger:         logger,
	}

	time.AfterFunc(time.Millisecond*200, func() {
		ctxCancel()
	})

	var cancelled bool

	go func() {
		rnr.Watch(ctx, ctxCancel, wg, false)
		cancelled = true
	}()

	<-time.After(time.Millisecond * 500)
	assert.True(t, cancelled)

	scriptsMgr.AssertCalled(t, "Get")
	scriptsMgr.AssertExpectations(t)

	assert.Equal(t, 1, logs.FilterMessage("error get scripts").FilterField(zap.Error(e)).Len())
}
