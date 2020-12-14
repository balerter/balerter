package alert

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshalLua(t *testing.T) {
	a := &Alert{
		name:       "foo",
		level:      LevelSuccess,
		lastChange: time.Date(2020, 01, 01, 01, 01, 01, 00, time.UTC),
		start:      time.Date(2020, 01, 01, 01, 01, 01, 00, time.UTC),
		count:      10,
	}

	res := a.MarshalLua()

	assert.Equal(t, "foo", res.RawGetString("name").String())
	assert.Equal(t, "success", res.RawGetString("level").String())
	assert.Equal(t, "1577840461", res.RawGetString("last_change").String())
	assert.Equal(t, "10", res.RawGetString("count").String())
}
