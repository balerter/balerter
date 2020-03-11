package alerts

import (
	"github.com/balerter/balerter/internal/alert/alert"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResource(t *testing.T) {

	data := []*alertManager.AlertInfo{
		{Level: alert.LevelError, Name: "foo"},
		{Level: alert.LevelSuccess, Name: "bar"},
	}

	r := newResource(data)

	assert.Equal(t, 2, len(r.items))

	rw := &http.TestResponseWriter{}

	err := r.render(rw)
	require.NoError(t, err)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Equal(t, `[{"name":"foo","level":"error"},{"name":"bar","level":"success"}]`, rw.Output)
}
