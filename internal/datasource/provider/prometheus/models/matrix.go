package models

import (
	"strings"
)

type Matrix []*SampleStream

func (m Matrix) String() string {
	var s []string

	for _, v := range m {
		s = append(s, v.String())
	}

	return strings.Join(s, "\n")
}

func (m Matrix) Type() string {
	return ValMatrix
}
