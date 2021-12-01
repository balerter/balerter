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

package stream

import (
	"context"
	"io"

	"go.uber.org/thriftrw/wire"
)

// BodyReader represents a type that can be read out from a stream.Reader.
type BodyReader interface {
	Decode(Reader) error
}

// Enveloper is the interface implemented by a type that can be written with
// an envelope via a stream.Writer.
type Enveloper interface {
	MethodName() string
	EnvelopeType() wire.EnvelopeType
	Encode(Writer) error
}

// RequestReader captures how to read from a request in a streaming fashion.
type RequestReader interface {
	// ReadRequest reads off the request envelope (if present) from an io.Reader,
	// using the provided BodyReader to read off the full request struct,
	// asserting the EnvelopeType (either OneWay or Unary) if an envlope exists.
	// A ResponseWriter that understands the enveloping used and the request's
	// body are returned.
	//
	// This allows a Thrift request handler to transparently read requests
	// regardless of whether the caller is configured to submit envelopes.
	ReadRequest(context.Context, wire.EnvelopeType, io.Reader, BodyReader) (ResponseWriter, error)
}

// ResponseWriter captures how to respond to a request in a streaming fashion.
type ResponseWriter interface {
	// WriteResponse writes a response to the Writer with the envelope
	// style of the corresponding request, and returns a stream.Writer to write
	// remaining un-enveloped response bytes. Once writing of the response is complete,
	// whether successful or not (error), users must call Close() on the stream.Writer.
	//
	// The EnvelopeType should be either wire.Reply or wire.Exception.
	WriteResponse(wire.EnvelopeType, io.Writer, Enveloper) error
}
