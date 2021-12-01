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

package wire

import "fmt"

// EnvelopeType is the type of data inside of the envelope.
type EnvelopeType int8

// List of envelope types, same as TMessageType in Apache Thrift.
const (
	Call      EnvelopeType = 1
	Reply     EnvelopeType = 2
	Exception EnvelopeType = 3
	OneWay    EnvelopeType = 4
)

// Envelope represents an enveloped value which includes metadata about
// the method, the type of data in the envelope, and the value.
type Envelope struct {
	Name  string
	Type  EnvelopeType
	SeqID int32
	Value Value
}

func (et EnvelopeType) String() string {
	switch et {
	case Call:
		return "Call"
	case Reply:
		return "Reply"
	case Exception:
		return "Exception"
	case OneWay:
		return "OneWay"
	default:
		return fmt.Sprintf("EnvelopeType(%v)", int8(et))
	}
}
