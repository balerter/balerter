package folder

import (
	"github.com/balerter/balerter/internal/config/scripts/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	p := New(file.File{
		Name:     "foo",
		Filename: "bar",
	})
	assert.Equal(t, "file.foo", p.name)
	assert.Equal(t, "bar", p.filename)
}

func TestProvider_Get_no_file(t *testing.T) {
	p := &Provider{filename: "/tmp/no-file"}
	_, err := p.Get()
	require.Error(t, err)
	assert.Equal(t, "open /tmp/no-file: no such file or directory", err.Error())
}

func TestProvider_Get_error_parse_meta(t *testing.T) {
	f, err := os.CreateTemp("", "")
	require.NoError(t, err)
	_, err = f.Write([]byte("-- @timeout adc\n\nfoo"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	p := &Provider{filename: f.Name()}
	_, err = p.Get()
	require.Error(t, err)
	assert.Equal(t, "error parse 'adc' to time duration, time: invalid duration \"adc\"", err.Error())
}

func TestProvider_Get(t *testing.T) {
	f, err := os.CreateTemp("", "")
	require.NoError(t, err)
	_, err = f.Write([]byte("foo"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	p := &Provider{filename: f.Name()}
	s, err := p.Get()
	require.NoError(t, err)
	require.Equal(t, 1, len(s))
	assert.Equal(t, "foo", string(s[0].Body))
}
