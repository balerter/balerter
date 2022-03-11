package models

import "time"

type Entry struct {
	Timestamp time.Time
	Line      string
}

type LabelSet map[string]string

type Stream struct {
	Labels  LabelSet `json:"stream"`
	Entries []Entry  `json:"values"`
}

type Streams []Stream

func (s Streams) Type() ResultType {
	return ResultTypeStream
}
