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
	"fmt"

	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

// WriteEnvelopeBegin writes the start of a strict envelope (contains an envelope version).
func (sw *StreamWriter) WriteEnvelopeBegin(eh stream.EnvelopeHeader) error {
	version := uint32(version1) | uint32(eh.Type)

	if err := sw.WriteInt32(int32(version)); err != nil {
		return err
	}

	if err := sw.WriteString(eh.Name); err != nil {
		return err
	}

	return sw.WriteInt32(eh.SeqID)
}

// WriteEnvelopeEnd writes the "end" of an envelope. Since there is no ending
// to an envelope, this is a no-op.
func (sw *StreamWriter) WriteEnvelopeEnd() error {
	return nil
}

// WriteLegacyEnvelopeBegin writes the start of a non-strict envelope (lacks an envelope version).
func (sw *StreamWriter) WriteLegacyEnvelopeBegin(eh stream.EnvelopeHeader) error {
	if err := sw.WriteString(eh.Name); err != nil {
		return err
	}

	if err := sw.writeByte(uint8(eh.Type)); err != nil {
		return err
	}

	return sw.WriteInt32(eh.SeqID)
}

// WriteLegacyEnvelopeEnd writes the "end" of a legacy envelope. Since there is
// no ending to a legacy envelope, this is a no-op.
func (sw *StreamWriter) WriteLegacyEnvelopeEnd() error {
	return nil
}

// ReadEnvelopeBegin reads the start of an Apache Thrift envelope. Thrift supports
// two kinds of envelopes: strict, and non-strict. See ReadEnveloped method
// for more information on enveloping.
func (sw *StreamReader) ReadEnvelopeBegin() (stream.EnvelopeHeader, error) {
	var eh stream.EnvelopeHeader

	val, err := sw.ReadInt32()
	if err != nil {
		return eh, err
	}

	if val > 0 {
		if eh, err = sw.readNonStrictEnvelope(val); err != nil {
			return eh, err
		}
	} else {
		if eh, err = sw.readStrictEnvelope(val); err != nil {
			return eh, err
		}
	}

	seqID, err := sw.ReadInt32()
	if err != nil {
		return eh, err
	}

	eh.SeqID = seqID
	return eh, nil
}

// readNonStrictEnvelope reads off a non-strict envelope as described by protocol.EnvelopeAgnosticProtocol.
func (sw *StreamReader) readNonStrictEnvelope(length int32) (stream.EnvelopeHeader, error) {
	var eh stream.EnvelopeHeader

	buf := make([]byte, length)
	for i := int32(0); i < length; i++ {
		i8, err := sw.ReadInt8()
		if err != nil {
			return eh, err
		}
		buf[i] = byte(i8)
	}

	typ, err := sw.ReadInt8()
	if err != nil {
		return eh, err
	}

	eh.Name = string(buf)
	eh.Type = wire.EnvelopeType(typ)
	return eh, nil
}

func (sw *StreamReader) readStrictEnvelope(ver int32) (stream.EnvelopeHeader, error) {
	var eh stream.EnvelopeHeader

	if v := uint32(ver) & versionMask; v != version1 {
		return eh, fmt.Errorf("cannot decode envelope of version: %v", v)
	}

	name, err := sw.ReadString()
	if err != nil {
		return eh, err
	}

	// Casting automatically truncates to the lowest 8 bits.
	eh.Type = wire.EnvelopeType(ver)
	eh.Name = name
	return eh, nil
}

// ReadEnvelopeEnd reads the "end" of an envelope.  Since there is no real
// envelope end, this is a no-op.
func (sw *StreamReader) ReadEnvelopeEnd() error {
	return nil
}
