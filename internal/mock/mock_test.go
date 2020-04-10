package mock

import (
	"github.com/balerter/balerter/internal/mock/registry"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	m := New("name", []string{"foo"}, zap.NewNop())

	assert.IsType(t, &ModuleMock{}, m)
	assert.Equal(t, "name", m.name)
	assert.Equal(t, []string{"foo"}, m.methods)
	assert.Equal(t, "name", m.Name())

	assert.NoError(t, m.Stop())
}

func TestClean(t *testing.T) {
	m := &ModuleMock{
		errors:   []string{"foo", "bar"},
		registry: registry.New(),
	}
	m.Clean()

	assert.Equal(t, 0, len(m.errors))
}

func TestGetLoader(t *testing.T) {
	m := &ModuleMock{
		methods: []string{
			"m1",
			"m2",
		},
	}

	f := m.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("on")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("assertCalled")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("assertNotCalled")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("m1")))
	assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString("m2")))
}
