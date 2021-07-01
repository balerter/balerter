package folder

import (
	"github.com/balerter/balerter/internal/config/scripts/folder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	p := New(folder.Folder{
		Name: "a",
		Path: "b",
		Mask: "c",
	})
	assert.IsType(t, &Provider{}, p)
	assert.Equal(t, "folder.a", p.name)
	assert.Equal(t, "b", p.path)
	assert.Equal(t, "c", p.mask)
}

func TestNew_default_mask(t *testing.T) {
	p := New(folder.Folder{})
	assert.Equal(t, "*.lua", p.mask)
}

func TestProvider_Get_bad_pattern(t *testing.T) {
	p := &Provider{
		name: "",
		path: "",
		mask: "\\",
	}

	_, err := p.Get()
	require.Error(t, err)
	assert.Equal(t, "syntax error in pattern", err.Error())
}

func TestProvider_Get_no_files(t *testing.T) {
	p := &Provider{
		name: "",
		path: "",
		mask: "*.noext",
	}

	ss, err := p.Get()
	require.NoError(t, err)
	assert.Equal(t, 0, len(ss))
}

func TestProvider_Get(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	prefix := strconv.Itoa(rand.Intn(1e6)) + "-"

	f1, err := os.CreateTemp("", prefix)
	require.NoError(t, err)
	_, err = f1.Write([]byte("foo1"))
	require.NoError(t, err)
	err = f1.Close()
	require.NoError(t, err)

	f2, err := os.CreateTemp("", prefix)
	require.NoError(t, err)
	_, err = f2.Write([]byte("foo2"))
	require.NoError(t, err)
	err = f2.Close()
	require.NoError(t, err)

	p := &Provider{
		name: "",
		path: os.TempDir() + "/",
		mask: prefix + "*",
	}

	ss, err := p.Get()
	require.NoError(t, err)
	assert.Equal(t, 2, len(ss))

	for _, s := range ss {
		if string(s.Body) != "foo1" && string(s.Body) != "foo2" {
			t.Fatal("unexpected file body")
		}
	}
}
