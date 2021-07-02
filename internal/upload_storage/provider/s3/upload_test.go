package s3

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"testing"
)

func TestS3_getArgs_no_upload_data(t *testing.T) {
	p := &Provider{}

	luaState := lua.NewState()

	_, _, err := p.getArgs(luaState)
	require.Error(t, err)
	assert.Equal(t, "upload data must be a string", err.Error())
}

func TestS3_getArgs_wrong_upload_data(t *testing.T) {
	p := &Provider{}

	luaState := lua.NewState()
	luaState.Push(lua.LNumber(42))

	_, _, err := p.getArgs(luaState)
	require.Error(t, err)
	assert.Equal(t, "upload data must be a string", err.Error())
}

func TestS3_getArgs_wrong_filename(t *testing.T) {
	p := &Provider{}

	luaState := lua.NewState()
	luaState.Push(lua.LString("foo"))
	luaState.Push(lua.LNumber(42))

	data, filename, err := p.getArgs(luaState)
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), data)
	assert.Equal(t, "", filename)
}

func TestS3_getArgs_without_filename(t *testing.T) {
	p := &Provider{}

	luaState := lua.NewState()
	luaState.Push(lua.LString("foo"))

	data, filename, err := p.getArgs(luaState)
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), data)
	assert.Equal(t, 30, len(filename))
}

func TestS3_getArgs_filename(t *testing.T) {
	p := &Provider{}

	luaState := lua.NewState()
	luaState.Push(lua.LString("foo"))
	luaState.Push(lua.LString("bar"))

	data, filename, err := p.getArgs(luaState)
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), data)
	assert.Equal(t, "bar", filename)
}

func TestS3_getArgs_filename_trim_suffix(t *testing.T) {
	p := &Provider{}

	luaState := lua.NewState()
	luaState.Push(lua.LString("foo"))
	luaState.Push(lua.LString("bar.png"))

	data, filename, err := p.getArgs(luaState)
	require.NoError(t, err)
	assert.Equal(t, []byte("foo"), data)
	assert.Equal(t, "bar", filename)
}

func TestProvider_upload_error_get_args(t *testing.T) {
	p := &Provider{
		name:     "",
		region:   "",
		endpoint: "",
		key:      "",
		secret:   "",
		bucket:   "",
		logger:   zap.NewNop(),
	}

	luaState := lua.NewState()

	n := p.uploadPNG(luaState)
	assert.Equal(t, 2, n)

	assert.Equal(t, "wrong arguments: upload data must be a string", luaState.Get(2).String())
}

func TestProvider_upload_error_upload(t *testing.T) {
	p := &Provider{
		name:     "",
		region:   "",
		endpoint: "",
		key:      "",
		secret:   "",
		bucket:   "",
		logger:   zap.NewNop(),
	}

	luaState := lua.NewState()
	luaState.Push(lua.LString("foo"))

	n := p.uploadPNG(luaState)
	assert.Equal(t, 2, n)

	assert.Equal(t, "error upload object: MissingRegion: could not find region configuration", luaState.Get(3).String())
}

//func TestProvider_upload(t *testing.T) {
//	p := &Provider{
//		name:     "",
//		region:   "local",
//		endpoint: "http://127.0.0.1:9000",
//		key:      "key",
//		secret:   "secret",
//		bucket:   "bucket",
//		logger:   zap.NewNop(),
//	}
//
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("foo"))
//
//	n := p.uploadPNG(luaState)
//	assert.Equal(t, 1, n)
//
//	assert.Equal(t, "", luaState.Get(2).String())
//	assert.Equal(t, "", luaState.Get(3).String())
//}
