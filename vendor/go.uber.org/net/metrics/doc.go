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

// Package metrics is a telemetry client designed for Uber's software
// networking team. It prioritizes performance on the hot path and integration
// with both push- and pull-based collection systems. Like Prometheus and
// Tally, it supports metrics tagged with arbitrary key-value pairs.
//
// Metric Names and Uniqueness
//
// Like Prometheus, but unlike Tally, metric names should be relatively long and
// descriptive - generally speaking, metrics from the same process shouldn't
// share names. (See the documentation for the Root struct below for a longer
// explanation of the uniqueness rules.) For example, prefer
// "grpc_successes_by_procedure" over "successes", since "successes" is common
// and vague. Where relevant, metric names should indicate their unit of
// measurement (e.g., "grpc_success_latency_ms").
//
// Counters and Gauges
//
// Counters represent monotonically increasing values, like a car's odometer.
// Gauges represent point-in-time readings, like a car's speedometer. Both
// counters and gauges expose not only write operations (set, add, increment,
// etc.), but also atomic reads. This makes them easy to integrate directly
// into your business logic: you can use them anywhere you'd otherwise use a
// 64-bit atomic integer.
//
// Histograms
//
// This package doesn't support analogs of Tally's timer or Prometheus's
// summary, because they can't be accurately aggregated at query time.
// Instead, it approximates distributions of values with histograms. These
// require more up-front work to set up, but are typically more accurate and
// flexible when queried. See https://prometheus.io/docs/practices/histograms/
// for a more detailed discussion of the trade-offs involved.
//
// Vectors
//
// Plain counters, gauges, and histograms have a fixed set of tags. However,
// it's common to encounter situations where a subset of a metric's tags vary
// constantly. For example, you might want to track the latency of your
// database queries by table: you know the database cluster, application name,
// and hostname at process startup, but you need to specify the table name
// with each query. To model these situations, this package uses vectors.
//
// Each vector is a local cache of metrics, so accessing them is quite fast.
// Within a vector, all metrics share a common set of constant tags and a list
// of variable tags. In our database query example, the constant tags are
// cluster, application, and hostname, and the only variable tag is table
// name. Usage examples are included in the documentation for each vector
// type.
//
// Push and Pull
//
// This package integrates with StatsD- and M3-based collection systems by
// periodically pushing differential updates. (Users can integrate with other
// push-based systems by implementing the push.Target interface.) It
// integrates with pull-based collectors by exposing an HTTP handler that
// supports Prometheus's text and protocol buffer exposition formats. Examples
// of both push and pull integration are included in the documentation for the
// root struct's Push and ServeHTTP methods.
//
// See Also
//
// If you're unfamiliar with Tally and Prometheus, you may want to consult
// their documentation:
//   https://godoc.org/github.com/uber-go/tally
//   https://godoc.org/github.com/prometheus/client_golang/prometheus
package metrics // import "go.uber.org/net/metrics"
