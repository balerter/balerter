package prometheus_models

import "strings"

type Streams []*SampleStream

func (m Streams) String() string {
	var s []string

	for _, v := range m {
		s = append(s, v.String())
	}

	return strings.Join(s, "\n")
}

func (m Streams) Type() string {
	return ValStreams
}
