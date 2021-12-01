// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package metrics

import (
	"sort"
	"time"
)

// A Snapshot is a point-in-time view of the state of any non-histogram
// metric.
type Snapshot struct {
	Name  string
	Tags  Tags
	Value int64
}

func (s Snapshot) less(other Snapshot) bool {
	if s.Name != other.Name {
		return s.Name < other.Name
	}
	return s.Tags.less(other.Tags)
}

// A HistogramSnapshot is a point-in-time view of the state of a Histogram.
type HistogramSnapshot struct {
	Name   string
	Tags   Tags
	Unit   time.Duration
	Values []int64 // rounded up to bucket upper bounds
}

func (l HistogramSnapshot) less(other HistogramSnapshot) bool {
	if l.Name != other.Name {
		return l.Name < other.Name
	}
	return l.Tags.less(other.Tags)
}

// A RootSnapshot exposes all the metrics contained in a Root and all its
// Scopes. It's useful in tests, but relatively expensive to construct.
type RootSnapshot struct {
	Counters   []Snapshot
	Gauges     []Snapshot
	Histograms []HistogramSnapshot
}

func (s *RootSnapshot) sort() {
	sort.Slice(s.Counters, func(i, j int) bool {
		return s.Counters[i].less(s.Counters[j])
	})
	sort.Slice(s.Gauges, func(i, j int) bool {
		return s.Gauges[i].less(s.Gauges[j])
	})
	sort.Slice(s.Histograms, func(i, j int) bool {
		return s.Histograms[i].less(s.Histograms[j])
	})
}

func (s *RootSnapshot) add(m metric) {
	switch v := m.(type) {
	case *Counter:
		s.Counters = append(s.Counters, v.snapshot())
	case *Gauge:
		s.Gauges = append(s.Gauges, v.snapshot())
	case *Histogram:
		s.Histograms = append(s.Histograms, v.snapshot())
	case *CounterVector:
		s.Counters = append(s.Counters, v.snapshot()...)
	case *GaugeVector:
		s.Gauges = append(s.Gauges, v.snapshot()...)
	case *HistogramVector:
		s.Histograms = append(s.Histograms, v.snapshot()...)
	}
}
