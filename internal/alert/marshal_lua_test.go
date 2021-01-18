package alert

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshalLua(t *testing.T) {
	a := &Alert{
		Name:       "foo",
		Level:      LevelSuccess,
		LastChange: time.Date(2020, 01, 01, 01, 01, 01, 00, time.UTC),
		Start:      time.Date(2020, 01, 01, 01, 01, 01, 00, time.UTC),
		Count:      10,
	}

	res := a.MarshalLua()

	assert.Equal(t, "foo", res.RawGetString("Name").String())
	assert.Equal(t, "success", res.RawGetString("Level").String())
	assert.Equal(t, "1577840461", res.RawGetString("last_change").String())
	assert.Equal(t, "10", res.RawGetString("Count").String())
}
