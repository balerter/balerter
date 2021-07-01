package kv

import (
	"fmt"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestModuleName(t *testing.T) {
	assert.Equal(t, "kv", ModuleName())
}

func TestMethods(t *testing.T) {
	m := Methods()
	assert.Equal(t, []string{
		"all",
		"get",
		"put",
		"delete",
		"upsert",
	}, m)
}

func TestNew(t *testing.T) {
	k := New(nil)
	assert.IsType(t, &KV{}, k)
}

func TestKV_GetLoader(t *testing.T) {
	kv := &KV{}

	f := kv.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}

func TestKV_Stop(t *testing.T) {
	kv := &KV{}
	assert.NoError(t, kv.Stop())
}

func Test_all_error(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("All").Return(map[string]string{}, fmt.Errorf("err1"))

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.all(ls)
	assert.Equal(t, 2, n)

	e := ls.Get(2)
	assert.Equal(t, "err1", e.String())
}

func Test_all(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("All").Return(map[string]string{"a": "b"}, nil)

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.all(ls)
	assert.Equal(t, 2, n)

	e := ls.Get(1)
	require.Equal(t, lua.LTTable, e.Type())

	assert.Equal(t, "b", e.(*lua.LTable).RawGetString("a").String())
}

func Test_get_error(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Get", mock.Anything).Return("", fmt.Errorf("err1"))

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.get(ls)
	assert.Equal(t, 2, n)

	e := ls.Get(2)
	assert.Equal(t, "err1", e.String())
}

func Test_get(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Get", mock.Anything).Return("foo", nil)

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.get(ls)
	assert.Equal(t, 2, n)

	e := ls.Get(1)
	assert.Equal(t, "foo", e.String())
}

func Test_put_error(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Put", mock.Anything, mock.Anything).Return(fmt.Errorf("err1"))

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.put(ls)
	assert.Equal(t, 1, n)

	e := ls.Get(1)
	assert.Equal(t, "err1", e.String())
}

func Test_put(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Put", mock.Anything, mock.Anything).Return(nil)

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.put(ls)
	assert.Equal(t, 0, n)
}

func Test_upsert_error(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Upsert", mock.Anything, mock.Anything).Return(fmt.Errorf("err1"))

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.upsert(ls)
	assert.Equal(t, 1, n)

	e := ls.Get(1)
	assert.Equal(t, "err1", e.String())
}

func Test_upsert(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Upsert", mock.Anything, mock.Anything).Return(nil)

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.upsert(ls)
	assert.Equal(t, 0, n)
}

func Test_delete_error(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Delete", mock.Anything).Return(fmt.Errorf("err1"))

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.delete(ls)
	assert.Equal(t, 1, n)

	e := ls.Get(1)
	assert.Equal(t, "err1", e.String())
}

func Test_delete(t *testing.T) {
	s := &coreStorage.KVMock{}
	s.On("Delete", mock.Anything).Return(nil)

	kv := &KV{engine: s}

	ls := lua.NewState()

	n := kv.delete(ls)
	assert.Equal(t, 0, n)
}

func TestKV_Name(t *testing.T) {
	kv := &KV{}
	assert.Equal(t, "kv", kv.Name())
}
