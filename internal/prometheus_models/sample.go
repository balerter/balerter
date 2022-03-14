package prometheus_models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type SampleValue float64

func (v SampleValue) String() string {
	return strconv.FormatFloat(float64(v), 'b', 16, 8)
}

func (v *SampleValue) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	vv, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*v = SampleValue(vv)
	return nil
}

type Metric map[string]string

type SamplePair struct {
	Timestamp Time
	Value     SampleValue
}

func (s *SamplePair) UnmarshalJSON(b []byte) error {
	v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
	return json.Unmarshal(b, &v)
}

func (s SamplePair) String() string {
	return fmt.Sprintf("%d => %f", s.Timestamp, s.Value)
}

type SampleMatrix struct {
	Metric Metric       `json:"metric"`
	Values []SamplePair `json:"values"`
}

func (s *SampleMatrix) String() string {
	var res []string
	m, _ := json.Marshal(s.Metric)
	res = append(res, string(m))
	for _, v := range s.Values {
		res = append(res, v.String())
	}
	return strings.Join(res, "\n")
}

type Sample struct {
	Metric Metric     `json:"metric"`
	Value  SamplePair `json:"value"`
}

func (s *Sample) String() string {
	var res []string
	m, _ := json.Marshal(s.Metric)
	res = append(res, string(m))
	res = append(res, s.Value.String())
	return strings.Join(res, "\n")
}

type StringValue string

func (s *StringValue) UnmarshalJSON(b []byte) error {
	*s = StringValue(strings.Trim(string(b), `"`))
	return nil
}

type StreamSamplePair struct {
	Timestamp Time
	Value     StringValue
}

func (s *StreamSamplePair) UnmarshalJSON(b []byte) error {
	v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
	return json.Unmarshal(b, &v)
}

func (s StreamSamplePair) String() string {
	return fmt.Sprintf("%d => %s", s.Timestamp, s.Value)
}

type SampleStream struct {
	Metric Metric             `json:"stream"`
	Values []StreamSamplePair `json:"values"`
}

func (s SampleStream) String() string {
	var res []string
	m, _ := json.Marshal(s.Metric)
	res = append(res, string(m))
	for _, v := range s.Values {
		res = append(res, v.String())
	}
	return strings.Join(res, "\n")
}
