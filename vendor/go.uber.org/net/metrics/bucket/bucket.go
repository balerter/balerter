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

// Package bucket provides utility functions for constructing and merging
// histogram buckets. Buckets are meant to be used in the base metrics
// package's HistogramSpec struct.
package bucket // import "go.uber.org/net/metrics/bucket"

import (
	"errors"
)

// NewRPCLatency returns a hand-crafted set of buckets useful for tracking the
// latency of RPCs (in milliseconds). Buckets range from 1 to 10000 (i.e.,
// 1ms-10s), getting less granular as latency increases.
func NewRPCLatency() []int64 {
	return []int64{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		12, 14, 16, 18, 20,
		25, 30, 35, 40, 45, 50,
		60, 70, 80, 90, 100,
		120, 140, 160, 180, 200,
		250, 300, 350, 400, 450, 500,
		600, 700, 800, 900, 1000,
		1500, 2000, 2500, 3000,
		4000, 5000, 7500, 10000,
	}
}

// NewExponential creates n exponential buckets, starting with the supplied
// initial value and increasing by a user-defined factor each time.
//
// If n is less than one, the initial value is less than one, or the factor
// is less than two, NewExponential returns a nil slice.
func NewExponential(initial, factor int64, n int) []int64 {
	if n < 1 || initial < 1 || factor < 2 {
		return nil
	}
	buckets := make([]int64, n)
	buckets[0] = initial
	for i := 1; i < len(buckets); i++ {
		buckets[i] = buckets[i-1] * factor
	}
	return buckets
}

// NewLinear creates n linear buckets, starting with the supplied
// initial value and increasing by a user-defined width.
//
// If n or width are less than one, NewLinear returns a nil slice.
func NewLinear(initial, width int64, n int) []int64 {
	if n < 1 || width < 1 {
		return nil
	}
	buckets := make([]int64, n)
	buckets[0] = initial
	for i := 1; i < len(buckets); i++ {
		buckets[i] = buckets[i-1] + width
	}
	return buckets
}

// Flatten concatenates multiple sets of buckets into a single slice. After
// flattening, it checks that the result is sorted and that no buckets are
// duplicated.
func Flatten(buckets ...[]int64) ([]int64, error) {
	merged := flatten(buckets)
	if asc := isAscending(merged); !asc {
		return nil, errors.New("after flattening, buckets are not strictly ascending")
	}
	return merged, nil
}

func isAscending(ints []int64) bool {
	// Don't use sort.IsSorted, since it permits duplicate elements.
	if len(ints) < 2 {
		return true
	}
	prev := ints[0]
	for j := 1; j < len(ints); j++ {
		if prev >= ints[j] {
			return false
		}
		prev = ints[j]
	}
	return true
}

func flatten(iss [][]int64) []int64 {
	var n int
	for _, is := range iss {
		n += len(is)
	}
	flat := make([]int64, 0, n)
	for _, is := range iss {
		flat = append(flat, is...)
	}
	return flat
}
