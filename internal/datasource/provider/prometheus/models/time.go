package models

import (
	"strconv"
)

type Time int64

func (t *Time) UnmarshalJSON(b []byte) error {
	v, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Time(v)

	return nil
}
