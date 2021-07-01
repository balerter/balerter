package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New(zap.NewNop())
	assert.IsType(t, &Service{}, s)
}

func Test_livenessHandler(t *testing.T) {
	s := &Service{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://domain.com", nil)

	s.livenessHandler(w, r)

	assert.Equal(t, "ok", w.Body.String())
}

type mockListener struct {
	mock.Mock
}

func (m *mockListener) Accept() (net.Conn, error) {
	a := m.Called()
	c := a.Get(0)
	if c == nil {
		return nil, a.Error(1)
	}
	return a.Get(0).(net.Conn), a.Error(1)
}

func (m *mockListener) Close() error {
	a := m.Called()
	return a.Error(0)
}

func (m *mockListener) Addr() net.Addr {
	a := m.Called()
	c := a.Get(0)
	if c == nil {
		return nil
	}
	return a.Get(0).(net.Addr)
}

func TestService_Run(t *testing.T) {
	srvMock := &http.Server{
		ReadTimeout:  time.Millisecond * 50,
		WriteTimeout: time.Millisecond * 50,
	}
	core, logs := observer.New(zap.ErrorLevel)

	s := &Service{
		logger: zap.New(core),
		server: srvMock,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	ln := &mockListener{}
	ln.On("Addr").Return(&net.IPAddr{})
	ln.On("Accept").Return(nil, fmt.Errorf("err1"))
	ln.On("Close").Return(fmt.Errorf("err2"))

	s.Run(ctx, cancel, wg, ln)

	select {
	case <-ctx.Done():
	case <-time.After(time.Millisecond * 500):
		t.Fatal("too long pause")
		return
	}

	wg.Wait()

	assert.Equal(t, 1, logs.FilterMessage("error serve service server").Len())
}
