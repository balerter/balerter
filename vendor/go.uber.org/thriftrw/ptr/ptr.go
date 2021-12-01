// Copyright (c) 2021 Uber Technologies, Inc.
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

// Package ptr provides helpers to convert basic types to pointers.
package ptr

// Bool converts a bool to a pointer
func Bool(x bool) *bool {
	return &x
}

// Int8 converts an int8 to a pointer
func Int8(x int8) *int8 {
	return &x
}

// Int16 converts an int32 to a pointer
func Int16(x int16) *int16 {
	return &x
}

// Int32 converts an int32 to a pointer
func Int32(x int32) *int32 {
	return &x
}

// Int64 converts an int64 to a pointer
func Int64(x int64) *int64 {
	return &x
}

// Float64 converts a float64 to a pointer
func Float64(x float64) *float64 {
	return &x
}

// String converts a string to a pointer
func String(x string) *string {
	return &x
}
