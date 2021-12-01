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

package envelope

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/internal/envelope/exception"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/wire"
)

// Enveloper is the interface implemented by a type that can be written with
// an envelope.
type Enveloper interface {
	MethodName() string
	EnvelopeType() wire.EnvelopeType
	ToWire() (wire.Value, error)
}

// Write writes an Envelope to the given writer.
func Write(p protocol.Protocol, w io.Writer, seqID int32, e Enveloper) error {
	body, err := e.ToWire()
	if err != nil {
		return err
	}
	return p.EncodeEnveloped(wire.Envelope{
		SeqID: seqID,
		Name:  e.MethodName(),
		Type:  e.EnvelopeType(),
		Value: body,
	}, w)
}

// ReadReply reads enveloped responses from the given reader.
func ReadReply(p protocol.Protocol, r io.ReaderAt) (_ wire.Value, seqID int32, _ error) {
	envelope, err := p.DecodeEnveloped(r)
	if err != nil {
		return wire.Value{}, 0, err
	}

	switch {
	case envelope.Type == wire.Reply:
		return envelope.Value, envelope.SeqID, nil
	case envelope.Type != wire.Exception:
		return envelope.Value, envelope.SeqID, fmt.Errorf("unknown envelope type for reply, got %v", envelope.Type)
	}

	// Decode the exception payload.
	ex := &exception.TApplicationException{}
	if err := ex.FromWire(envelope.Value); err != nil {
		return envelope.Value, envelope.SeqID, fmt.Errorf("failed to decode exception: %v", err)
	}

	return envelope.Value, envelope.SeqID, ex
}
