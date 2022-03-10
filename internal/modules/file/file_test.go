package file

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestModuleName(t *testing.T) {
	assert.Equal(t, "file", ModuleName())
}

func TestName(t *testing.T) {
	f := &File{}
	assert.Equal(t, "file", f.Name())
}

func TestMethods(t *testing.T) {
	assert.Equal(t, []string{"save", "load"}, Methods())
}

func TestFile_Stop(t *testing.T) {
	f := &File{}
	assert.NoError(t, f.Stop())
}

func TestFile_GetLoader(t *testing.T) {
	h := &File{}

	loader := h.GetLoader(nil)

	L := lua.NewState()
	n := loader(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)
	assert.Equal(t, lua.LTFunction, v.RawGetString("save").Type())
	assert.Equal(t, lua.LTFunction, v.RawGetString("load").Type())
}

func TestFile_save(t *testing.T) {
	filename := filepath.Join(os.TempDir(), strconv.Itoa(rand.Int())+".txt")

	payload := "bar"

	l := lua.NewState()
	l.Push(lua.LString(filename))
	l.Push(lua.LString(payload))

	f := &File{}

	n := f.save(l)
	assert.Equal(t, 0, n)

	data, errRead := os.ReadFile(filename)
	if errRead != nil {
		t.Fatal(errRead)
	}

	assert.Equal(t, string(data), payload)
}

func TestFile_load(t *testing.T) {
	filename := filepath.Join(os.TempDir(), strconv.Itoa(rand.Int())+".txt")
	payload := "bar"

	errWrite := os.WriteFile(filename, []byte(payload), 0644)
	if errWrite != nil {
		t.Fatal(errWrite)
	}

	l := lua.NewState()
	l.Push(lua.LString(filename))

	f := &File{}

	n := f.load(l)
	assert.Equal(t, 1, n)

	v := l.Get(2).String()

	assert.Equal(t, v, payload)
}
