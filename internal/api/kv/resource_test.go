package kv

import (
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResource(t *testing.T) {
	data := map[string]string{
		"f1": "v1",
		"f2": "v2",
	}

	r := newResource(data)

	assert.Equal(t, 2, len(r.items))

	rw := &httpTestify.TestResponseWriter{}

	err := r.render(rw)
	require.NoError(t, err)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Equal(t, `[{"name":"f1","value":"v1"},{"name":"f2","value":"v2"}]`, rw.Output)

}
