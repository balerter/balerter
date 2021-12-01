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
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	promproto "github.com/prometheus/client_model/go"
	"go.uber.org/net/metrics/push"
)

const _defaultCollectionSize = 128

// A core is a collection of metrics. Uniqueness is enforced with two checks,
// explained in the documentation of the Root struct. The checks are slightly
// stricter than those used by the official Prometheus client.
//
// The test suite for metric uniqueness is well-commented and explores the
// consequences of these two rules.
type core struct {
	sync.RWMutex

	dimsByName map[string]string
	ids        map[string]struct{}
	metrics    []metric
	gatherer   prometheus.Gatherer
}

func newCore() *core {
	c := &core{
		dimsByName: make(map[string]string, _defaultCollectionSize),
		ids:        make(map[string]struct{}, _defaultCollectionSize),
		metrics:    make([]metric, 0, _defaultCollectionSize),
	}
	c.gatherer = prometheus.GathererFunc(func() ([]*promproto.MetricFamily, error) {
		c.RLock()
		protos := make([]*promproto.MetricFamily, 0, len(c.metrics))
		for _, m := range c.metrics {
			p := m.proto()
			if p != nil && len(p.Metric) > 0 {
				protos = append(protos, p)
			}
		}
		c.RUnlock()
		return protos, nil
	})
	return c
}

func (c *core) register(m metric) error {
	id := newDigester()
	defer id.free()

	meta := m.describe()
	meta.writeID(id)

	c.Lock()
	if existing, ok := c.dimsByName[*meta.Name]; ok && existing != meta.Dims {
		c.Unlock()
		return fmt.Errorf("a metric with name %q and different tag "+
			"names is already registered", *meta.Name)
	}
	if _, ok := c.ids[string(id.digest())]; ok {
		c.Unlock()
		return fmt.Errorf("a metric with name %q and the same constant "+
			"tag names and values is already registered", *meta.Name)
	}
	c.dimsByName[*meta.Name] = meta.Dims
	c.ids[string(id.digest())] = struct{}{}
	c.metrics = append(c.metrics, m)
	c.Unlock()

	return nil
}

func (c *core) snapshot() *RootSnapshot {
	c.RLock()
	defer c.RUnlock()
	s := &RootSnapshot{}
	for _, m := range c.metrics {
		s.add(m)
	}
	s.sort()
	return s
}

func (c *core) push(target push.Target) {
	c.RLock()
	for _, m := range c.metrics {
		m.push(target)
	}
	c.RUnlock()
}
