package meta

import (
	"testing"
	"time"

	"github.com/balerter/balerter/internal/modules"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestNew(t *testing.T) {
	m := New(nil)

	assert.IsType(t, &Meta{}, m)
}

func TestName(t *testing.T) {
	m := &Meta{}

	assert.Equal(t, "meta", m.Name())
}

func TestGetLoader(t *testing.T) {
	m := &Meta{}

	j := &modules.JobMock{}

	f := m.GetLoader(j)

	L := lua.NewState()

	n := f(L)

	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.Equal(t, lua.LTFunction, v.RawGetString(method).Type())
	}
}

func TestStop(t *testing.T) {
	m := &Meta{}

	assert.NoError(t, m.Stop())
}

func Test_priorExecutionTime(t *testing.T) {
	m := &Meta{}

	j := &modules.JobMock{
		GetPriorExecutionTimeFunc: func() time.Duration {
			return time.Second
		},
	}

	L := lua.NewState()
	f := m.priorExecutionTime(j)
	n := f(L)
	assert.Equal(t, 1, n)
	v := L.Get(1)
	assert.Equal(t, lua.LTNumber, v.Type())
	assert.Equal(t, "1", v.String())
}
