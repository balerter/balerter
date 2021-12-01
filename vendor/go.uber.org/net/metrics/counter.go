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
	"fmt"
	"sort"

	promproto "github.com/prometheus/client_model/go"
	"go.uber.org/net/metrics/push"
)

// A Counter is a monotonically increasing value, like a car's odometer. All
// its exported methods are safe to use concurrently, and nil *Counters are
// safe no-op implementations.
type Counter struct {
	val    value
	pusher push.Counter
}

func newCounter(m metadata) *Counter {
	return &Counter{val: newValue(m)}
}

func newDynamicCounter(m metadata, variableTagPairs []string) metric {
	return &Counter{val: newDynamicValue(m, variableTagPairs)}
}

// Add increases the value of the counter and returns the new value. Since
// counters must be monotonically increasing, passing a negative number just
// returns the current value (without modifying it).
func (c *Counter) Add(n int64) int64 {
	if c == nil {
		return 0
	}
	if n <= 0 {
		return c.val.Load()
	}
	return c.val.Add(n)
}

// Inc increments the counter's value by one and returns the new value.
func (c *Counter) Inc() int64 {
	if c == nil {
		return 0
	}
	return c.val.Inc()
}

// Load returns the counter's current value.
func (c *Counter) Load() int64 {
	if c == nil {
		return 0
	}
	return c.val.Load()
}

func (c *Counter) describe() metadata {
	return c.val.meta
}

func (c *Counter) snapshot() Snapshot {
	return c.val.snapshot()
}

func (c *Counter) proto() *promproto.MetricFamily {
	return &promproto.MetricFamily{
		Name:   c.val.meta.Name,
		Help:   c.val.meta.Help,
		Type:   promproto.MetricType_COUNTER.Enum(),
		Metric: []*promproto.Metric{c.metric()},
	}
}

func (c *Counter) metric() *promproto.Metric {
	n := float64(c.val.Load())
	return &promproto.Metric{
		Label:   c.val.tagPairs,
		Counter: &promproto.Counter{Value: &n},
	}
}

func (c *Counter) push(target push.Target) {
	if c.val.meta.DisablePush {
		return
	}
	if c.pusher == nil {
		c.pusher = target.NewCounter(push.Spec{
			Name: *c.val.meta.Name,
			Tags: zip(c.val.tagPairs),
		})
	}
	c.pusher.Set(c.Load())
}

// A CounterVector is a collection of Counters that share a name and some
// constant tags, but also have a consistent set of variable tags. All
// exported methods are safe to use concurrently. Nil *CounterVectors are safe
// to use and always return no-op counters.
//
// For a general description of vector types, see the package-level
// documentation.
type CounterVector struct {
	vector
}

func newCounterVector(m metadata) *CounterVector {
	return &CounterVector{vector{
		meta:    m,
		factory: newDynamicCounter,
		metrics: make(map[string]metric, _defaultCollectionSize),
	}}
}

// Get retrieves the counter with the supplied variable tag names and values
// from the vector, creating one if necessary. The variable tags must be
// supplied in the same order used when creating the vector.
//
// Get returns an error if the number or order of tags is incorrect.
func (cv *CounterVector) Get(variableTagPairs ...string) (*Counter, error) {
	if cv == nil {
		return nil, nil
	}
	m, err := cv.getOrCreate(variableTagPairs)
	if err != nil {
		return nil, err
	}
	return m.(*Counter), nil
}

// MustGet behaves exactly like Get, but panics on errors. If code using this
// method is covered by unit tests, this is safe.
func (cv *CounterVector) MustGet(variableTagPairs ...string) *Counter {
	if cv == nil {
		return nil
	}
	c, err := cv.Get(variableTagPairs...)
	if err != nil {
		panic(fmt.Sprintf("failed to get counter: %v", err))
	}
	return c
}

func (cv *CounterVector) describe() metadata {
	return cv.meta
}

func (cv *CounterVector) proto() *promproto.MetricFamily {
	mf := &promproto.MetricFamily{
		Name: cv.meta.Name,
		Help: cv.meta.Help,
		Type: promproto.MetricType_COUNTER.Enum(),
	}
	cv.metricsMu.RLock()
	protos := make([]*promproto.Metric, 0, len(cv.metrics))
	for _, metric := range cv.metrics {
		protos = append(protos, metric.(*Counter).metric())
	}
	cv.metricsMu.RUnlock()
	sort.Slice(protos, func(i, j int) bool {
		return protos[i].String() < protos[j].String()
	})
	mf.Metric = protos
	return mf
}
