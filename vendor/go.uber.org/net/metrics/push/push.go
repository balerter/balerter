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

// Package push integrates go.uber.org/net/metrics with push-based telemetry
// systems like Graphite and M3. To integrate with such a system, implement
// the Target interface and use the Push method on metrics.Root.
//
// See the go.uber.org/net/metrics/tallypush package for an example
// integration with both StatsD- and M3-based systems.
package push // import "go.uber.org/net/metrics/push"

// A Target bridges the metrics package's representations of counters, gauges,
// and histograms with push-based telemetry systems. Targets are designed to
// work with the metrics.Root struct's Push method, so they don't need to be
// safe for concurrent use.
//
// A concrete implementation of this interface that works with StatsD and M3
// is available in the go.uber.org/net/metrics/tallypush package.
type Target interface {
	NewCounter(Spec) Counter
	NewGauge(Spec) Gauge
	NewHistogram(HistogramSpec) Histogram
}

// A Spec configures counters and gauges.
type Spec struct {
	Name string
	Tags map[string]string
}

// A HistogramSpec configures histograms.
type HistogramSpec struct {
	Spec

	Buckets []int64 // upper bounds, inclusive
}

// A Counter models monotonically increasing values, like a car's odometer.
// Implementations should expect to be called with the total accumulated value
// of the counter.
//
// Implementations do not need to be safe for concurrent use.
type Counter interface {
	Set(total int64)
}

// A Gauge models moment-in-time measurements, like a car's speedometer.
// Implementations should expect to be called with the current value of the
// gauge.
//
// Implementations do not need to be safe for concurrent use.
type Gauge interface {
	Set(total int64)
}

// A Histogram approximates a distribution of values. Implementations should
// expect to be called with the upper bound of a bucket and the total
// accumulated number of observations in that bucket.
//
// Implementations do not need to be safe for concurrent use.
type Histogram interface {
	Set(bucket int64, total int64)
}
