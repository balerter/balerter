package runner

import (
	"bytes"
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

type scriptManagerMock struct {
	mock.Mock
}

func (m *scriptManagerMock) Get() ([]*script.Script, error) {
	a := m.Called()
	v := a.Get(0)
	if v == nil {
		return nil, a.Error(1)
	}
	return v.([]*script.Script), a.Error(1)
}

type storageManagerMock struct {
	mock.Mock
}

func (m *storageManagerMock) Get() []modules.Module {
	a := m.Called()
	v := a.Get(0)
	if v == nil {
		return nil
	}
	return v.([]modules.Module)
}

func TestRunner_RunScript_error_get_script(t *testing.T) {
	m := &scriptManagerMock{}
	m.On("Get").Return(nil, fmt.Errorf("err1"))

	rnr := &Runner{
		scriptsManager: m,
	}

	req, err := http.NewRequest("POST", "localhost", nil)
	require.NoError(t, err)

	err = rnr.RunScript("foo", req)
	require.Error(t, err)
	assert.Equal(t, "err1", err.Error())
}

type badReader struct {
	mock.Mock
}

func (m *badReader) Read(b []byte) (int, error) {
	a := m.Called(b)
	return a.Int(0), a.Error(1)
}

func TestRunner_RunScript_error_create_luaState(t *testing.T) {
	s := &storageManagerMock{}
	s.On("Get").Return(nil)
	d := &storageManagerMock{}
	d.On("Get").Return(nil)
	m := &scriptManagerMock{}
	m.On("Get").Return([]*script.Script{{
		Name:       "foo",
		Body:       nil,
		CronValue:  "",
		Timeout:    0,
		Ignore:     false,
		Channels:   nil,
		IsTest:     false,
		TestTarget: "",
	}}, nil)

	rnr := &Runner{
		scriptsManager:  m,
		storagesManager: s,
		dsManager:       d,
		logger:          zap.NewNop(),
		jobs:            make(chan *Job, 1),
	}

	r := &badReader{}
	r.On("Read", mock.Anything).Return(0, fmt.Errorf("err1"))

	req, err := http.NewRequest("POST", "localhost", r)
	require.NoError(t, err)

	err = rnr.RunScript("foo", req)
	require.Error(t, err)
	assert.Equal(t, "error init api module, err1", err.Error())
}

func TestRunner_RunScript_script_not_found(t *testing.T) {
	m := &scriptManagerMock{}
	m.On("Get").Return(nil, nil)

	rnr := &Runner{
		scriptsManager: m,
		logger:         zap.NewNop(),
		jobs:           make(chan *Job, 1),
	}

	req, err := http.NewRequest("POST", "localhost", nil)
	require.NoError(t, err)

	err = rnr.RunScript("foo", req)
	require.Error(t, err)
	assert.Equal(t, "script foo not found", err.Error())
}

func TestRunner_RunScript(t *testing.T) {
	s := &storageManagerMock{}
	s.On("Get").Return(nil)
	d := &storageManagerMock{}
	d.On("Get").Return(nil)
	m := &scriptManagerMock{}
	m.On("Get").Return([]*script.Script{{
		Name:       "foo",
		Body:       nil,
		CronValue:  "",
		Timeout:    0,
		Ignore:     false,
		Channels:   nil,
		IsTest:     false,
		TestTarget: "",
	}}, nil)

	rnr := &Runner{
		scriptsManager:  m,
		storagesManager: s,
		dsManager:       d,
		logger:          zap.NewNop(),
		jobs:            make(chan *Job, 1),
	}

	req, err := http.NewRequest("POST", "localhost", bytes.NewBuffer([]byte("bar")))
	require.NoError(t, err)

	err = rnr.RunScript("foo", req)
	require.NoError(t, err)

	var j *Job

	select {
	case j = <-rnr.jobs:
	default:
		t.Fatalf("error read job from the channel")
	}

	assert.Equal(t, "foo", j.name)
}
