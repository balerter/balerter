package prometheus_models

import (
	"strconv"
	"strings"
)

type Time int64

func (t *Time) UnmarshalJSON(b []byte) error {
	v, err := strconv.Atoi(strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}

	*t = Time(v)

	return nil
}
