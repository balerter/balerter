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
	"errors"
	"fmt"
	"math"
	"time"
)

// A Spec configures Counters, Gauges, CounterVectors, and GaugeVectors.
type Spec struct {
	Name        string   // required: metric name, should be fairly long and descriptive
	Help        string   // required: displayed on HTTP pages
	ConstTags   Tags     // optional: constant tags
	VarTags     []string // variable tags, required for vectors and forbidden otherwise
	DisablePush bool     // reduces load on system we're pushing to (if any)
}

func (s Spec) validate() error {
	if len(s.Name) == 0 {
		return errors.New("all metrics must have a name")
	}
	if s.Help == "" {
		return errors.New("metric help must not be empty")
	}
	return nil
}

func (s Spec) validateScalar() error {
	if err := s.validate(); err != nil {
		return err
	}
	if len(s.VarTags) > 0 {
		return errors.New("only vectors may have variable tags")
	}
	return nil
}

func (s Spec) validateVector() error {
	if err := s.validate(); err != nil {
		return err
	}
	if len(s.VarTags) == 0 {
		return errors.New("vectors must have variable tags")
	}
	return nil
}

// A HistogramSpec configures Histograms and HistogramVectors.
type HistogramSpec struct {
	Spec

	// Durations are exposed as simple numbers, not strings or rich objects.
	// Unit specifies the desired granularity for histogram observations. For
	// example, an observation of time.Second with a unit of time.Millisecond is
	// exposed as 1000. Typically, the unit should also be part of the metric
	// name.
	Unit time.Duration
	// Upper bounds (inclusive) for the histogram buckets in terms of the unit.
	// A catch-all bucket for large observations is automatically created, if
	// necessary.
	Buckets []int64
}

func (hs HistogramSpec) validateScalar() error {
	if err := hs.validateHistogram(); err != nil {
		return err
	}
	return hs.Spec.validateScalar()
}

func (hs HistogramSpec) validateVector() error {
	if err := hs.validateHistogram(); err != nil {
		return err
	}
	return hs.Spec.validateVector()
}

func (hs HistogramSpec) validateHistogram() error {
	if hs.Unit < 1 {
		return fmt.Errorf("duration unit must be positive, got %v", hs.Unit)
	}
	if len(hs.Buckets) == 0 {
		return fmt.Errorf("must specify some buckets")
	}
	prev := int64(math.MinInt64)
	for _, upper := range hs.Buckets {
		if upper <= prev {
			return fmt.Errorf("bucket upper bounds must be sorted in increasing order")
		}
		prev = upper
	}
	return nil
}
