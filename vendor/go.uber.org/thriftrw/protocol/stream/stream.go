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

// Package stream provides streaming implementations of encoding and decoding
// Thrift values.
package stream

import (
	"io"

	"go.uber.org/thriftrw/wire"
)

// Protocol defines a specific way for a Thrift value to be encoded or
// decoded, implemented in a streaming fashion.
type Protocol interface {
	// Writer returns a streaming implementation of an encoder for a
	// Thrift value.
	Writer(w io.Writer) Writer

	// Reader returns a streaming implementation of a decoder for a
	// Thrift value.
	Reader(r io.Reader) Reader
}

// EnvelopeHeader represents the envelope of a response or a request which includes
// metadata about the method, the type of data in the envelope, and the value.
// It is equivalent of `wire.Envelope`, but for streaming purposes.
type EnvelopeHeader struct {
	Name  string
	Type  wire.EnvelopeType
	SeqID int32
}

// FieldHeader defines the metadata needed to define the beginning of a field
// in a Thrift value.
type FieldHeader struct {
	ID   int16
	Type wire.Type
}

// MapHeader defines the metadata needed to define the beginning of a map in a
// Thrift value.
type MapHeader struct {
	KeyType   wire.Type
	ValueType wire.Type
	Length    int
}

// SetHeader defines the metadata needed to define the beginning of a set in a
// Thrift value.
type SetHeader struct {
	Length int
	Type   wire.Type
}

// ListHeader defines the metadata needed to define the beginning of a list in a
// Thrift value.
type ListHeader struct {
	Length int
	Type   wire.Type
}

// Writer defines an encoder for a Thrift value, implemented in a streaming
// fashion.
type Writer interface {
	WriteBool(b bool) error
	WriteInt8(i int8) error
	WriteInt16(i int16) error
	WriteInt32(i int32) error
	WriteInt64(i int64) error
	WriteString(s string) error
	WriteDouble(f float64) error
	WriteBinary(b []byte) error
	WriteStructBegin() error
	WriteStructEnd() error
	WriteFieldBegin(f FieldHeader) error
	WriteFieldEnd() error
	WriteMapBegin(m MapHeader) error
	WriteMapEnd() error
	WriteSetBegin(s SetHeader) error
	WriteSetEnd() error
	WriteListBegin(l ListHeader) error
	WriteListEnd() error

	WriteEnvelopeBegin(eh EnvelopeHeader) error
	WriteEnvelopeEnd() error

	Close() error
}

// Reader defines an decoder for a Thrift value, implemented in a streaming
// fashion.
type Reader interface {
	ReadBool() (bool, error)
	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadString() (string, error)
	ReadDouble() (float64, error)
	ReadBinary() ([]byte, error)
	ReadStructBegin() error
	ReadStructEnd() error
	ReadFieldBegin() (FieldHeader, bool, error)
	ReadFieldEnd() error
	ReadListBegin() (ListHeader, error)
	ReadListEnd() error
	ReadSetBegin() (SetHeader, error)
	ReadSetEnd() error
	ReadMapBegin() (MapHeader, error)
	ReadMapEnd() error
	Close() error

	ReadEnvelopeBegin() (EnvelopeHeader, error)
	ReadEnvelopeEnd() error

	// Skip skips over the bytes of the wire type and any applicable headers.
	Skip(w wire.Type) error
}
