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

package binary

import (
	"io"

	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

// Responder captures how to respond to a request, concerning whether and what
// kind of envelope to use, how to match the sequence identifier of the
// corresponding request.
type Responder interface {
	EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error
}

// noEnvelopeResponder responds to a request without an envelope.
type noEnvelopeResponder struct{}

var _ stream.ResponseWriter = &noEnvelopeResponder{}

func (noEnvelopeResponder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	return Default.Encode(v, w)
}

func (noEnvelopeResponder) WriteResponse(et wire.EnvelopeType, w io.Writer, ev stream.Enveloper) error {
	writer := NewStreamWriter(w)
	defer writer.Close()

	return ev.Encode(writer)
}

// NoEnvelopeResponder responds to a request without an envelope.
var NoEnvelopeResponder = &noEnvelopeResponder{}

// EnvelopeV0Responder responds to requests with a non-strict (unversioned) envelope.
type EnvelopeV0Responder struct {
	Name  string
	SeqID int32
}

var _ stream.ResponseWriter = &EnvelopeV0Responder{}

// EncodeResponse writes the response to the writer using a non-strict
// envelope.
func (r EnvelopeV0Responder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	writer := BorrowWriter(w)
	err := writer.WriteLegacyEnveloped(wire.Envelope{
		Name:  r.Name,
		Type:  t,
		SeqID: r.SeqID,
		Value: v,
	})
	ReturnWriter(writer)
	return err
}

// WriteResponse writes an envelope to the writer (non-strict envelope) and
// returns a borrowed stream.Writer. Callers must call Close() on stream.Writer once finished.
func (r EnvelopeV0Responder) WriteResponse(et wire.EnvelopeType, w io.Writer, ev stream.Enveloper) error {
	writer := NewStreamWriter(w)
	defer writer.Close()

	if err := writer.WriteLegacyEnvelopeBegin(stream.EnvelopeHeader{
		Name:  r.Name,
		Type:  et,
		SeqID: r.SeqID,
	}); err != nil {
		return err
	}

	if err := ev.Encode(writer); err != nil {
		return err
	}

	return writer.WriteLegacyEnvelopeEnd()
}

// EnvelopeV1Responder responds to requests with a strict, version 1 envelope.
type EnvelopeV1Responder struct {
	Name  string
	SeqID int32
}

var _ stream.ResponseWriter = &EnvelopeV1Responder{}

// EncodeResponse writes the response to the writer using a strict, version 1
// envelope.
func (r EnvelopeV1Responder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	writer := BorrowWriter(w)
	err := writer.WriteEnveloped(wire.Envelope{
		Name:  r.Name,
		Type:  t,
		SeqID: r.SeqID,
		Value: v,
	})
	ReturnWriter(writer)
	return err
}

// WriteResponse writes an envelope to the writer (strict envelope) and returns a
// borrowed stream.Writer. Callers must call Close() on stream.Writer once finished.
func (r EnvelopeV1Responder) WriteResponse(et wire.EnvelopeType, w io.Writer, ev stream.Enveloper) error {
	writer := NewStreamWriter(w)
	defer writer.Close()

	if err := writer.WriteEnvelopeBegin(stream.EnvelopeHeader{
		Name:  r.Name,
		Type:  et,
		SeqID: r.SeqID,
	}); err != nil {
		return err
	}

	if err := ev.Encode(writer); err != nil {
		return err
	}

	return writer.WriteEnvelopeEnd()
}
