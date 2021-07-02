package s3

import (
	"github.com/balerter/balerter/internal/config/storages/upload/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestModuleName(t *testing.T) {
	assert.Equal(t, "s3.foo", ModuleName("foo"))
}

func TestMethods(t *testing.T) {
	assert.Equal(t, []string{"uploadPNG"}, Methods())
}

func TestNew(t *testing.T) {
	p, err := New(s3.S3{}, zap.NewNop())
	require.NoError(t, err)
	assert.IsType(t, &Provider{}, p)
}

func TestProvider_Name(t *testing.T) {
	p := &Provider{name: "foo"}
	assert.Equal(t, "foo", p.name)
}

func TestProvider_Stop(t *testing.T) {
	p := &Provider{}
	assert.NoError(t, p.Stop())
}

func TestProvider_GetLoader(t *testing.T) {
	p := &Provider{}

	f := p.GetLoader(nil)

	L := lua.NewState()
	n := f(L)
	assert.Equal(t, 1, n)

	v := L.Get(1).(*lua.LTable)

	for _, method := range Methods() {
		assert.IsType(t, &lua.LFunction{}, v.RawGet(lua.LString(method)))
	}
}
