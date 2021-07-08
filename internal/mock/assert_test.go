package mock

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func Test_assert_no_first(t *testing.T) {
	mm := &ModuleMock{
		logger: zap.NewNop(),
	}
	f := mm.assert(false)

	ls := lua.NewState()

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 1, len(mm.errors))
	assert.Equal(t, "query should have a method name as first argument", mm.errors[0])
}

func Test_assert_method_not_string(t *testing.T) {
	mm := &ModuleMock{
		logger: zap.NewNop(),
	}
	f := mm.assert(false)

	ls := lua.NewState()
	ls.Push(lua.LNumber(42))

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 1, len(mm.errors))
	assert.Equal(t, "method name must be a string", mm.errors[0])
}

func Test_assert_method_empty(t *testing.T) {
	mm := &ModuleMock{
		logger: zap.NewNop(),
	}
	f := mm.assert(false)

	ls := lua.NewState()
	ls.Push(lua.LString(" "))

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 1, len(mm.errors))
	assert.Equal(t, "method name must be a not empty", mm.errors[0])
}

func Test_assert_error_add_assert(t *testing.T) {
	r := &registryMock{}
	r.On("AddAssert", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("err1"))
	mm := &ModuleMock{
		logger:   zap.NewNop(),
		registry: r,
	}
	f := mm.assert(false)

	ls := lua.NewState()
	ls.Push(lua.LString("foo"))

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 1, len(mm.errors))
	assert.Equal(t, "error register assert: err1", mm.errors[0])
}

func Test_assert(t *testing.T) {
	r := &registryMock{}
	r.On("AddAssert", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mm := &ModuleMock{
		logger:   zap.NewNop(),
		registry: r,
	}
	f := mm.assert(false)

	ls := lua.NewState()
	ls.Push(lua.LString("foo"))

	n := f(ls)
	assert.Equal(t, 0, n)
	assert.Equal(t, 0, len(mm.errors))
}
