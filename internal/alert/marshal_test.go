package alert

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	etalon = []byte{0x3, 0x62, 0x61, 0x72, 0x2, 0x1, 0x0, 0x0, 0x0, 0xe, 0xd5, 0x9d, 0xe6, 0x4d, 0x0, 0x0, 0x0,
		0x1, 0xff, 0xff, 0x1, 0x0, 0x0, 0x0, 0xe, 0xd5, 0x9d, 0xe6, 0x4d, 0x0, 0x0, 0x0, 0x1, 0xff, 0xff, 0xa}
)

func TestAlert_Marshal(t *testing.T) {
	now := time.Date(2020, 01, 01, 01, 01, 01, 01, time.UTC)

	a := &Alert{
		Name:       "bar",
		Level:      2,
		LastChange: now,
		Start:      now,
		Count:      10,
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

	assert.Equal(t, "bar", a.Name)
	assert.Equal(t, Level(2), a.Level)
	assert.Equal(t, now, a.LastChange)
	assert.Equal(t, now, a.Start)
	assert.Equal(t, 10, a.Count)
}

func TestAlert_Unmarshal_Errors(t *testing.T) {
	a := &Alert{}
	var err error

	err = a.Unmarshal(nil)
	require.Error(t, err)
	assert.Equal(t, "error decode alert Name", err.Error())

	err = a.Unmarshal([]byte{0x9F})
	require.Error(t, err)
	assert.Equal(t, "error decode alert Name", err.Error())

	err = a.Unmarshal([]byte{0x03, 'b', 'a'})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x02})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())
	assert.Equal(t, "bac", a.Name)
	assert.Equal(t, Level(2), a.Level)

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())
	assert.Equal(t, "0001-01-01T00:00:00Z", a.LastChange.Format(time.RFC3339))

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	require.Error(t, err)
	assert.Equal(t, "source too small", err.Error())
	assert.Equal(t, "0001-01-01T00:00:00Z", a.LastChange.Format(time.RFC3339))
	assert.Equal(t, "0001-01-01T00:00:00Z", a.Start.Format(time.RFC3339))

	err = a.Unmarshal([]byte{0x03, 'b', 'a', 'c', 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05})
	require.NoError(t, err)
	assert.Equal(t, 5, a.Count)
}
