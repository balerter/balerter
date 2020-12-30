package alerts

import (
	"testing"
)

// TODO: wip

func TestResource(t *testing.T) {
	//var data []*alert2.Alert
	//
	//a1 := alert2.AcquireAlert()
	//a1.SetName("foo")
	//a1.UpdateLevel(alert2.LevelError)
	//a1.Inc()
	//a1.Inc()
	//updateAt1 := a1.GetLastChangeTime().Format(time.RFC3339)
	//data = append(data, a1)
	//
	//a2 := alert2.AcquireAlert()
	//a2.SetName("bar")
	//a2.UpdateLevel(alert2.LevelSuccess)
	//a2.Inc()
	//a2.Inc()
	//updateAt2 := a1.GetLastChangeTime().Format(time.RFC3339)
	//data = append(data, a2)
	//
	//r := newResource(data)
	//
	//assert.Equal(t, 2, len(r.items))
	//
	//rw := &httpTestify.TestResponseWriter{}
	//
	//err := r.render(rw)
	//require.NoError(t, err)
	//
	//assert.Equal(t, 200, rw.StatusCode)
	//assert.Contains(t, rw.Output, `{"name":"foo","level":"error","count":2,"updated_at":"`+updateAt1+`"}`)
	//assert.Contains(t, rw.Output, `{"name":"bar","level":"success","count":2,"updated_at":"`+updateAt2+`"}`)
}
