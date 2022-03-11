package models

import (
	"strings"
)

type Vector []*Sample

func (m Vector) String() string {
	var s []string

	for _, v := range m {
		s = append(s, v.String())
	}

	return strings.Join(s, "\n")
}

func (m Vector) Type() string {
	return ValVector
}
