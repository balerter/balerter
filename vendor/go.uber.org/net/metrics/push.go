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
	"time"

	"go.uber.org/net/metrics/push"
)

type pusher struct {
	core    *core
	stop    chan struct{}
	stopped chan struct{}
	ticker  *time.Ticker
	target  push.Target
}

func newPusher(c *core, target push.Target, tick time.Duration) *pusher {
	return &pusher{
		core:    c,
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
		ticker:  time.NewTicker(tick),
		target:  target,
	}
}

func (p *pusher) Start() {
	defer close(p.stopped)
	// When stopping, do one last export to catch any stragglers.
	defer p.core.push(p.target)

	for {
		select {
		case <-p.stop:
			return
		case <-p.ticker.C:
			p.core.push(p.target)
		}
	}
}

func (p *pusher) Stop() {
	p.ticker.Stop()
	close(p.stop)
	<-p.stopped
}
