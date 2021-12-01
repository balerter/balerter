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
	"math"
	"sort"
	"sync"
	"time"

	promproto "github.com/prometheus/client_model/go"
	"go.uber.org/atomic"
	"go.uber.org/net/metrics/push"
)

type bucket struct {
	atomic.Int64

	upper int64 // bucket upper bound, inclusive
}

type buckets []*bucket

func newBuckets(upperBounds []int64) buckets {
	bs := make(buckets, 0, len(upperBounds)+1)
	for _, upper := range upperBounds {
		bs = append(bs, &bucket{upper: upper})
	}
	if upperBounds[len(upperBounds)-1] != math.MaxInt64 {
		bs = append(bs, &bucket{upper: math.MaxInt64})
	}
	return bs
}

func (bs buckets) get(val int64) *bucket {
	// Binary search to find the correct bucket for this observation. Bucket
	// upper bounds are inclusive.
	i, j := 0, len(bs)
	for i < j {
		h := i + (j-i)/2
		if val > bs[h].upper {
			i = h + 1
		} else {
			j = h
		}
	}
	return bs[i]
}

// A Histogram approximates a distribution of values. They're both more
// efficient and easier to aggregate than Prometheus summaries or M3 timers.
// For a discussion of the tradeoffs between histograms and timers/summaries,
// see https://prometheus.io/docs/practices/histograms/.
//
// All exported methods are safe to use concurrently, and nil *Histograms are
// valid no-op implementations.
type Histogram struct {
	meta     metadata
	unit     time.Duration
	bounds   []int64
	buckets  buckets
	sum      atomic.Int64 // required by Prometheus
	pusher   push.Histogram
	tagPairs []*promproto.LabelPair
}

func newHistogram(m metadata, unit time.Duration, uppers []int64) *Histogram {
	return &Histogram{
		buckets:  newBuckets(uppers),
		meta:     m,
		unit:     unit,
		bounds:   uppers,
		tagPairs: m.MergeTags(nil /* variable tag vals */),
	}
}

// Observe finds the correct bucket for the supplied duration and increments
// its counter. This is purely a convenience - it's equivalent to dividing the
// duration by the histogram's unit and calling IncBucket directly.
func (h *Histogram) Observe(d time.Duration) {
	if h == nil {
		return
	}
	h.IncBucket(int64(d / h.unit))
}

// IncBucket bypasses the time-based Observe API and increments a histogram
// bucket directly. It finds the correct bucket for the supplied value and
// adds one to its counter.
func (h *Histogram) IncBucket(n int64) {
	if h == nil {
		return
	}
	bucket := h.buckets.get(n)
	bucket.Inc()
	h.sum.Add(n)
}

func (h *Histogram) describe() metadata {
	return h.meta
}

func (h *Histogram) snapshot() HistogramSnapshot {
	return HistogramSnapshot{
		Name:   *h.meta.Name,
		Tags:   zip(h.tagPairs),
		Unit:   h.unit,
		Values: h.observations(),
	}
}

func (h *Histogram) observations() []int64 {
	var obs []int64
	for _, b := range h.buckets {
		n := b.Load()
		for i := int64(0); i < n; i++ {
			obs = append(obs, b.upper)
		}
	}
	return obs
}

func (h *Histogram) proto() *promproto.MetricFamily {
	return &promproto.MetricFamily{
		Name:   h.meta.Name,
		Help:   h.meta.Help,
		Type:   promproto.MetricType_HISTOGRAM.Enum(),
		Metric: []*promproto.Metric{h.metric()},
	}
}

func (h *Histogram) metric() *promproto.Metric {
	n := uint64(0)
	promBuckets := make([]*promproto.Bucket, 0, len(h.buckets)-1)
	for _, b := range h.buckets {
		n += uint64(b.Load())
		if b.upper == math.MaxInt64 {
			// Prometheus doesn't want us to export the final catch-all bucket.
			continue
		}
		upper := float64(b.upper)
		promBuckets = append(promBuckets, &promproto.Bucket{
			CumulativeCount: &n,
			UpperBound:      &upper,
		})
	}

	sum := float64(h.sum.Load())
	return &promproto.Metric{
		Label: h.tagPairs,
		Histogram: &promproto.Histogram{
			SampleCount: &n,
			SampleSum:   &sum,
			Bucket:      promBuckets,
		},
	}
}

func (h *Histogram) push(target push.Target) {
	if h.meta.DisablePush {
		return
	}
	if h.pusher == nil {
		h.pusher = target.NewHistogram(push.HistogramSpec{
			Spec: push.Spec{
				Name: *h.meta.Name,
				Tags: zip(h.tagPairs),
			},
			Buckets: h.bounds,
		})
	}
	for _, bucket := range h.buckets {
		h.pusher.Set(bucket.upper, bucket.Load())
	}
}

// A HistogramVector is a collection of Histograms that share a name and some
// constant tags, but also have a consistent set of variable tags. All
// exported methods are safe to use concurrently. Nil *HistogramVectors are
// safe to use and always return no-op histograms.
//
// For a general description of vector types, see the package-level
// documentation.
type HistogramVector struct {
	meta   metadata
	unit   time.Duration
	bounds []int64

	histogramsMu sync.RWMutex
	histograms   map[string]*Histogram // key is variable tag vals
}

func newHistogramVector(m metadata, unit time.Duration, uppers []int64) *HistogramVector {
	return &HistogramVector{
		meta:       m,
		unit:       unit,
		bounds:     uppers,
		histograms: make(map[string]*Histogram, _defaultCollectionSize),
	}
}

// Get retrieves the histogram with the supplied variable tag names and values
// from the vector, creating one if necessary. The variable tags must be
// supplied in the same order used when creating the vector.
//
// Get returns an error if the number or order of tags is incorrect.
func (hv *HistogramVector) Get(variableTagPairs ...string) (*Histogram, error) {
	if hv == nil {
		return nil, nil
	}
	if err := hv.meta.ValidateVariableTags(variableTagPairs); err != nil {
		return nil, err
	}
	digester := newDigester()
	for i := 0; i < len(variableTagPairs)/2; i++ {
		digester.add("", scrubTagValue(variableTagPairs[i*2+1]))
	}

	hv.histogramsMu.RLock()
	h, ok := hv.histograms[string(digester.digest())]
	hv.histogramsMu.RUnlock()
	if ok {
		digester.free()
		return h, nil
	}

	hv.histogramsMu.Lock()
	h, err := hv.newHistogram(digester.digest(), variableTagPairs)
	hv.histogramsMu.Unlock()
	digester.free()

	return h, err
}

// MustGet behaves exactly like Get, but panics on errors. If code using this
// method is covered by unit tests, this is safe.
func (hv *HistogramVector) MustGet(variableTagPairs ...string) *Histogram {
	if hv == nil {
		return nil
	}
	h, err := hv.Get(variableTagPairs...)
	if err != nil {
		panic(fmt.Sprintf("failed to get histogram: %v", err))
	}
	return h
}

func (hv *HistogramVector) newHistogram(key []byte, variableTagPairs []string) (*Histogram, error) {
	h, ok := hv.histograms[string(key)]
	if ok {
		return h, nil
	}
	h = &Histogram{
		buckets:  newBuckets(hv.bounds),
		meta:     hv.meta,
		unit:     hv.unit,
		bounds:   hv.bounds,
		tagPairs: hv.meta.MergeTags(variableTagPairs),
	}
	hv.histograms[string(key)] = h
	return h, nil
}

func (hv *HistogramVector) describe() metadata {
	return hv.meta
}

func (hv *HistogramVector) snapshot() []HistogramSnapshot {
	hv.histogramsMu.RLock()
	defer hv.histogramsMu.RUnlock()
	snaps := make([]HistogramSnapshot, 0, len(hv.histograms))
	for _, h := range hv.histograms {
		snaps = append(snaps, h.snapshot())
	}
	return snaps
}

func (hv *HistogramVector) proto() *promproto.MetricFamily {
	hv.histogramsMu.RLock()
	protos := make([]*promproto.Metric, 0, len(hv.histograms))
	for _, h := range hv.histograms {
		protos = append(protos, h.metric())
	}
	hv.histogramsMu.RUnlock()
	sort.Slice(protos, func(i, j int) bool {
		return protos[i].String() < protos[j].String()
	})

	return &promproto.MetricFamily{
		Name:   hv.meta.Name,
		Help:   hv.meta.Help,
		Type:   promproto.MetricType_HISTOGRAM.Enum(),
		Metric: protos,
	}
}

func (hv *HistogramVector) push(target push.Target) {
	hv.histogramsMu.RLock()
	for _, m := range hv.histograms {
		m.push(target)
	}
	hv.histogramsMu.RUnlock()
}
