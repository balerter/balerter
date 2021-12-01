// Copyright (c) 2015 Uber Technologies, Inc.
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

// +build gofuzz

package binary

import (
	"bytes"
	"fmt"

	"go.uber.org/thriftrw/wire"
)

func Fuzz(data []byte) int {
	reader := NewReader(bytes.NewReader(data))
	value, pos, err := reader.ReadValue(wire.TStruct, 0)
	if err != nil || pos != int64(len(data)) {
		return 0
	}
	if err := wire.EvaluateValue(value); err != nil {
		return 0
	}

	buffer := bytes.Buffer{}
	writer := BorrowWriter(&buffer)
	if err := writer.WriteValue(value); err != nil {
		panic(fmt.Sprintf("error encoding %v: %v", value, err))
	}
	ReturnWriter(writer)

	if encoded := buffer.Bytes(); !bytes.Equal(data, encoded) {
		panic(fmt.Sprintf(
			"encoding mismatch for %v:\n\t   %#v (got)\n\t!= %#v (expected)\n",
			value, encoded, data,
		))
	}

	return 1
}
