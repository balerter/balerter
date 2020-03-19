package alerts

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestResource(t *testing.T) {
	var data []*alert.Alert

	a1 := alert.AcquireAlert()
	a1.SetName("foo")
	a1.UpdateLevel(alert.LevelError)
	a1.Inc()
	a1.Inc()
	updateAt1 := a1.GetLastChangeTime().Format(time.RFC3339)
	data = append(data, a1)

	a2 := alert.AcquireAlert()
	a2.SetName("bar")
	a2.UpdateLevel(alert.LevelSuccess)
	a2.Inc()
	a2.Inc()
	updateAt2 := a1.GetLastChangeTime().Format(time.RFC3339)
	data = append(data, a2)

	r := newResource(data)

	assert.Equal(t, 2, len(r.items))

	rw := &httpTestify.TestResponseWriter{}

	err := r.render(rw)
	require.NoError(t, err)

	assert.Equal(t, 200, rw.StatusCode)
	assert.Equal(t, `[{"name":"foo","level":"error","count":2,"updated_at":"`+updateAt1+`"},{"name":"bar","level":"success","count":2,"updated_at":"`+updateAt2+`"}]`, rw.Output)
}
