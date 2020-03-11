package alerts

import (
	"github.com/balerter/balerter/internal/alert/alert"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestResource(t *testing.T) {

	data := []*alertManager.AlertInfo{
		{Level: alert.LevelError, Name: "foo", Count: 5, LastChange: time.Date(2020, 01, 01, 10, 10, 10, 0, time.UTC)},
		{Level: alert.LevelSuccess, Name: "bar", Count: 5, LastChange: time.Date(2020, 02, 02, 12, 12, 12, 0, time.UTC)},
	}

	r := newResource(data)

	assert.Equal(t, 2, len(r.items))

	rw := &http.TestResponseWriter{}

	err := r.render(rw)
	require.NoError(t, err)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Equal(t, `[{"name":"foo","level":"error","count":5,"updated_at":"2020-01-01T10:10:10Z"},{"name":"bar","level":"success","count":5,"updated_at":"2020-02-02T12:12:12Z"}]`, rw.Output)
}
