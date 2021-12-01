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

const (
	versionMask = 0xffff0000
	version1    = 0x80010000
)

// WriteEnveloped writes enveloped value using the strict envelope.
func (bw *Writer) WriteEnveloped(e wire.Envelope) error {
	if err := bw.sw.WriteEnvelopeBegin(
		stream.EnvelopeHeader{
			Name:  e.Name,
			Type:  e.Type,
			SeqID: e.SeqID,
		},
	); err != nil {
		return err
	}

	if err := bw.WriteValue(e.Value); err != nil {
		return err
	}

	return bw.sw.WriteEnvelopeEnd()
}

// WriteLegacyEnveloped writes enveloped value using the non-strict envelope
// (non-strict lacks an envelope version).
func (bw *Writer) WriteLegacyEnveloped(e wire.Envelope) error {
	if err := bw.sw.WriteLegacyEnvelopeBegin(
		stream.EnvelopeHeader{
			Name:  e.Name,
			Type:  e.Type,
			SeqID: e.SeqID,
		},
	); err != nil {
		return err
	}

	if err := bw.WriteValue(e.Value); err != nil {
		return err
	}

	return bw.sw.WriteLegacyEnvelopeEnd()
}

// ReadEnveloped reads an Apache Thrift envelope
//
// Thrift supports two kinds of envelopes: strict, and non-strict.
//
// Non-strict envelopes:
// Name (4 byte length prefixed string)
// Type ID (1 byte)
// Sequence ID (4 bytes)
//
// Strict envelopes:
//
// Version | Type ID (4 bytes)
// Name (4 byte length prefixed string)
// Sequence ID (4 bytes)
//
// When reading payloads, we need to support both strict and non-strict
// payloads. To do this, we read the first 4 byte. Non-strict payloads
// will always have a size >= 0, while strict payloads have selected
// version numbers such that the value will always be negative.
func (bw *Reader) ReadEnveloped() (wire.Envelope, error) {
	var e wire.Envelope
	val, off, err := bw.ReadValue(wire.TI32, 0)
	if err != nil {
		return wire.Envelope{}, err
	}

	initial := val.GetI32()
	if initial > 0 {
		e, off, err = bw.readNonStrictNameType()
	} else {
		e, off, err = bw.readStrictNameType(initial, off)
	}
	if err != nil {
		return e, err
	}

	val, off, err = bw.ReadValue(wire.TI32, off)
	e.SeqID = val.GetI32()
	if err != nil {
		return e, err
	}

	e.Value, _, err = bw.ReadValue(wire.TStruct, off)
	if err != nil {
		return wire.Envelope{}, err
	}

	return e, nil
}

func (bw *Reader) readStrictNameType(initial int32, off int64) (wire.Envelope, int64, error) {
	var e wire.Envelope

	if v := uint32(initial) & versionMask; v != version1 {
		return e, off, fmt.Errorf("cannot decode envelope of version: %v", v)
	}

	// This will truncate the bits that are not required.
	e.Type = wire.EnvelopeType(initial)

	name, off, err := bw.ReadValue(wire.TBinary, off)
	e.Name = name.GetString()
	return e, off, err
}

func (bw *Reader) readNonStrictNameType() (wire.Envelope, int64, error) {
	var e wire.Envelope

	name, off, err := bw.ReadValue(wire.TBinary, 0)
	if err != nil {
		return e, off, err
	}
	e.Name = name.GetString()

	val, off, err := bw.ReadValue(wire.TI8, off)
	if err != nil {
		return e, off, err
	}
	e.Type = wire.EnvelopeType(val.GetI8())

	return e, off, nil
}
