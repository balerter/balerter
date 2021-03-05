package manager

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	cm := corestorage.NewMock("")

	a := New("", cm, cm, nil, nil)
	assert.IsType(t, &API{}, a)
}

type httpServerMock struct {
	mock.Mock
}

func (m *httpServerMock) Serve(l net.Listener) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *httpServerMock) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestRun_error_serve(t *testing.T) {
	m := &httpServerMock{}

	core, logs := observer.New(zap.DebugLevel)

	a := &API{
		server: m,
		logger: zap.New(core),
	}

	m.On("Serve", mock.Anything).Return(fmt.Errorf("err1"))
	m.On("Shutdown", mock.Anything).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	a.Run(ctx, cancel, wg, nil)

	assert.Equal(t, 1, logs.FilterMessage("error serve api server").Len())
}

func TestRun_error_shutdown(t *testing.T) {
	m := &httpServerMock{}

	core, logs := observer.New(zap.DebugLevel)

	a := &API{
		server: m,
		logger: zap.New(core),
	}

	m.On("Serve", mock.Anything).Return(nil)
	m.On("Shutdown", mock.Anything).Return(fmt.Errorf("err1"))

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	a.Run(ctx, cancel, wg, nil)

	assert.Equal(t, 1, logs.FilterMessage("error shutdown api server").Len())
}
