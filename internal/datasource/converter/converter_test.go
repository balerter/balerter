package converter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertFunctions(t *testing.T) {
	v := 1.12
	assert.Equal(t, "1.12", FromFloat64(&v))

	d := time.Date(2020, 2, 3, 4, 5, 6, 7, time.UTC)
	assert.Equal(t, "2020-02-03", FromDate(&d))

	d = time.Date(2020, 2, 3, 4, 5, 6, 7, time.UTC)
	assert.Equal(t, "2020-02-03T04:05:06Z", FromDateTime(&d))

	s := "value"
	assert.Equal(t, "value", FromString(&s))

	b := true
	assert.Equal(t, "true", FromBoolean(&b))

	b = false
	assert.Equal(t, "false", FromBoolean(&b))

	ui := uint(100)
	assert.Equal(t, "100", FromUInt(&ui))

	i := 100
	assert.Equal(t, "100", FromInt(&i))

	i = -100
	assert.Equal(t, "-100", FromInt(&i))
}
