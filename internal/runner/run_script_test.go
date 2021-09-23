package runner

import (
	"bytes"
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestRunner_RunScript_error_get_script(t *testing.T) {
	m := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return nil, fmt.Errorf("err1")
		},
	}

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
}

func (m *badReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("err1")
}

func TestRunner_RunScript_error_create_luaState(t *testing.T) {
	s := &storagesManagerMock{
		GetFunc: func() []modules.Module {
			return nil
		},
	}
	d := &storagesManagerMock{
		GetFunc: func() []modules.Module {
			return nil
		},
	}
	m := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return []*script.Script{{
				Name:       "foo",
				Body:       nil,
				CronValue:  "",
				Timeout:    0,
				Ignore:     false,
				Channels:   nil,
				IsTest:     false,
				TestTarget: "",
			}}, nil
		},
	}

	rnr := &Runner{
		scriptsManager:  m,
		storagesManager: s,
		dsManager:       d,
		logger:          zap.NewNop(),
		jobs:            make(chan job, 1),
		cron:            cron.New(),
	}

	r := &badReader{}

	req, err := http.NewRequest("POST", "localhost", r)
	require.NoError(t, err)

	err = rnr.RunScript("foo", req)
	require.Error(t, err)
	assert.Equal(t, "error init api module, err1", err.Error())
}

func TestRunner_RunScript_script_not_found(t *testing.T) {
	m := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return nil, nil
		},
	}

	rnr := &Runner{
		scriptsManager: m,
		logger:         zap.NewNop(),
		jobs:           make(chan job, 1),
	}

	req, err := http.NewRequest("POST", "localhost", nil)
	require.NoError(t, err)

	err = rnr.RunScript("foo", req)
	require.Error(t, err)
	assert.Equal(t, "script foo not found", err.Error())
}

func TestRunner_RunScript(t *testing.T) {
	s := &storagesManagerMock{
		GetFunc: func() []modules.Module {
			return nil
		},
	}
	d := &storagesManagerMock{
		GetFunc: func() []modules.Module {
			return nil
		},
	}
	m := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return []*script.Script{{
				Name:       "foo",
				Body:       nil,
				CronValue:  "",
				Timeout:    0,
				Ignore:     false,
				Channels:   nil,
				IsTest:     false,
				TestTarget: "",
			}}, nil
		},
	}

	rnr := &Runner{
		scriptsManager:  m,
		storagesManager: s,
		dsManager:       d,
		logger:          zap.NewNop(),
		jobs:            make(chan job, 1),
		cron:            cron.New(),
	}

	req, err := http.NewRequest("POST", "localhost", bytes.NewBuffer([]byte("bar")))
	require.NoError(t, err)

	err = rnr.RunScript("foo", req)
	require.NoError(t, err)

	var j job

	select {
	case j = <-rnr.jobs:
	default:
		t.Fatalf("error read job from the channel")
	}

	assert.Equal(t, "foo", j.Name())
}
