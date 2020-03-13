package alert

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	etalon = []byte{0x3, 0x62, 0x61, 0x72, 0x2, 0x1, 0x0, 0x0, 0x0, 0xe, 0xd5, 0x9d, 0xe6, 0x4d, 0x0, 0x0, 0x0, 0x1, 0xff, 0xff, 0x1, 0x0, 0x0, 0x0, 0xe, 0xd5, 0x9d, 0xe6, 0x4d, 0x0, 0x0, 0x0, 0x1, 0xff, 0xff, 0xa}
)

func TestAlert_Marshal(t *testing.T) {
	now := time.Date(2020, 01, 01, 01, 01, 01, 01, time.UTC)

	a := &Alert{
		name:       "bar",
		level:      2,
		lastChange: now,
		start:      now,
		count:      10,
	}

	res, err := a.Marshal()
	require.NoError(t, err)

	assert.Equal(t, etalon, res)
}

func TestAlert_Unmarshal(t *testing.T) {
	now := time.Date(2020, 01, 01, 01, 01, 01, 01, time.UTC)

	a := &Alert{}

	err := a.Unmarshal(etalon)
	require.NoError(t, err)

	assert.Equal(t, "bar", a.name)
	assert.Equal(t, Level(2), a.level)
	assert.Equal(t, now, a.lastChange)
	assert.Equal(t, now, a.start)
	assert.Equal(t, 10, a.count)
}

func TestAlert_Unmarshal_Errors(t *testing.T) {
	a := &Alert{}
	var err error

	err = a.Unmarshal(nil)
	require.Error(t, err)
	assert.Equal(t, "error decode alert name", err.Error())

	err = a.Unmarshal([]byte{0x9F})
	require.Error(t, err)
	assert.Equal(t, "error decode alert name", err.Error())

	err = a.Unmarshal([]byte{0x03, 'b', 'a'})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x02})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())
	assert.Equal(t, "bac", a.name)
	assert.Equal(t, Level(2), a.level)

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())
	assert.Equal(t, "0001-01-01T00:00:00Z", a.lastChange.Format(time.RFC3339))

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())
	assert.Equal(t, "0001-01-01T00:00:00Z", a.lastChange.Format(time.RFC3339))
	assert.Equal(t, "0001-01-01T00:00:00Z", a.start.Format(time.RFC3339))

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05})
	require.NoError(t, err)
	assert.Equal(t, 5, a.count)
}
