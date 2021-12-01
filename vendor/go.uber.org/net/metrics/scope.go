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

// A Scope is a collection of tagged metrics.
type Scope struct {
	core      *core
	constTags Tags
}

func newScope(c *core, tags Tags) *Scope {
	return &Scope{
		core:      c,
		constTags: tags,
	}
}

// Tagged creates a new scope with new constant tags merged into the existing
// tags (if any). Tag names and values are automatically scrubbed, with
// invalid characters replaced by underscores.
func (s *Scope) Tagged(tags Tags) *Scope {
	if s == nil {
		return nil
	}
	newTags := make(Tags, len(s.constTags)+len(tags))
	for k, v := range s.constTags {
		newTags[k] = v
	}
	for k, v := range tags {
		newTags[scrubName(k)] = scrubTagValue(v)
	}
	return newScope(s.core, newTags)
}

// Counter constructs a new Counter.
func (s *Scope) Counter(spec Spec) (*Counter, error) {
	if s == nil {
		return nil, nil
	}
	spec = s.addConstTags(spec)
	if err := spec.validateScalar(); err != nil {
		return nil, err
	}
	meta, err := newMetadata(spec)
	if err != nil {
		return nil, err
	}
	c := newCounter(meta)
	if err := s.core.register(c); err != nil {
		return nil, err
	}
	return c, nil
}

// Gauge constructs a new Gauge.
func (s *Scope) Gauge(spec Spec) (*Gauge, error) {
	if s == nil {
		return nil, nil
	}
	spec = s.addConstTags(spec)
	if err := spec.validateScalar(); err != nil {
		return nil, err
	}
	meta, err := newMetadata(spec)
	if err != nil {
		return nil, err
	}
	g := newGauge(meta)
	if err := s.core.register(g); err != nil {
		return nil, err
	}
	return g, nil
}

// Histogram constructs a new Histogram.
func (s *Scope) Histogram(spec HistogramSpec) (*Histogram, error) {
	if s == nil {
		return nil, nil
	}
	spec.Spec = s.addConstTags(spec.Spec)
	if err := spec.validateScalar(); err != nil {
		return nil, err
	}
	meta, err := newMetadata(spec.Spec)
	if err != nil {
		return nil, err
	}
	h := newHistogram(meta, spec.Unit, spec.Buckets)
	if err := s.core.register(h); err != nil {
		return nil, err
	}
	return h, nil
}

// CounterVector constructs a new CounterVector.
func (s *Scope) CounterVector(spec Spec) (*CounterVector, error) {
	if s == nil {
		return nil, nil
	}
	spec = s.addConstTags(spec)
	if err := spec.validateVector(); err != nil {
		return nil, err
	}
	meta, err := newMetadata(spec)
	if err != nil {
		return nil, err
	}
	cv := newCounterVector(meta)
	if err := s.core.register(cv); err != nil {
		return nil, err
	}
	return cv, nil
}

// GaugeVector constructs a new GaugeVector.
func (s *Scope) GaugeVector(spec Spec) (*GaugeVector, error) {
	if s == nil {
		return nil, nil
	}
	spec = s.addConstTags(spec)
	if err := spec.validateVector(); err != nil {
		return nil, err
	}
	meta, err := newMetadata(spec)
	if err != nil {
		return nil, err
	}
	gv := newGaugeVector(meta)
	if err := s.core.register(gv); err != nil {
		return nil, err
	}
	return gv, nil
}

// HistogramVector constructs a new HistogramVector.
func (s *Scope) HistogramVector(spec HistogramSpec) (*HistogramVector, error) {
	if s == nil {
		return nil, nil
	}
	spec.Spec = s.addConstTags(spec.Spec)
	if err := spec.validateVector(); err != nil {
		return nil, err
	}
	meta, err := newMetadata(spec.Spec)
	if err != nil {
		return nil, err
	}
	hv := newHistogramVector(meta, spec.Unit, spec.Buckets)
	if err := s.core.register(hv); err != nil {
		return nil, err
	}
	return hv, nil
}

func (s *Scope) addConstTags(spec Spec) Spec {
	if len(s.constTags) == 0 {
		return spec
	}
	tags := make(Tags, len(s.constTags)+len(spec.ConstTags))
	for k, v := range s.constTags {
		tags[k] = v
	}
	for k, v := range spec.ConstTags {
		tags[k] = v
	}
	spec.ConstTags = tags
	return spec
}
