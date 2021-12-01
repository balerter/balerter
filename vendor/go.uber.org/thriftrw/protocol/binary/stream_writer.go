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
	"math"
	"sync"

	"go.uber.org/thriftrw/protocol/stream"
)

var streamWriterPool = sync.Pool{
	New: func() interface{} {
		return &StreamWriter{}
	}}

// StreamWriter implements basic logic for writing the Thrift Binary Protocol
// to an io.Writer.
type StreamWriter struct {
	writer io.Writer

	// This buffer is re-used every time we need a slice of up to 8 bytes.
	buffer [8]byte
}

// NewStreamWriter fetches a StreamWriter from the system that will write
// its output to the given io.Writer.
//
// This StreamWriter must be returned back using ReturnStreamWriter.
func NewStreamWriter(w io.Writer) *StreamWriter {
	streamWriter := streamWriterPool.Get().(*StreamWriter)
	streamWriter.writer = w
	return streamWriter
}

// returnStreamWriter returns a previously borrowed StreamWriter back to the
// system.
func returnStreamWriter(sw *StreamWriter) {
	sw.writer = nil
	streamWriterPool.Put(sw)
}

func (sw *StreamWriter) write(bs []byte) error {
	_, err := sw.writer.Write(bs)
	return err
}

func (sw *StreamWriter) writeByte(b byte) error {
	bs := sw.buffer[0:1]
	bs[0] = b
	return sw.write(bs)
}

// WriteBool encodes a boolean
func (sw *StreamWriter) WriteBool(b bool) error {
	if b {
		return sw.writeByte(1)
	}
	return sw.writeByte(0)
}

// WriteInt8 encodes an int8
func (sw *StreamWriter) WriteInt8(i int8) error {
	return sw.writeByte(byte(i))
}

// WriteInt16 encodes an int16
func (sw *StreamWriter) WriteInt16(i int16) error {
	bs := sw.buffer[0:2]
	bigEndian.PutUint16(bs, uint16(i))
	return sw.write(bs)
}

// WriteInt32 encodes an int32
func (sw *StreamWriter) WriteInt32(i int32) error {
	bs := sw.buffer[0:4]
	bigEndian.PutUint32(bs, uint32(i))
	return sw.write(bs)
}

// WriteInt64 encodes an int64
func (sw *StreamWriter) WriteInt64(i int64) error {
	bs := sw.buffer[0:8]
	bigEndian.PutUint64(bs, uint64(i))
	return sw.write(bs)
}

// WriteString encodes a string
func (sw *StreamWriter) WriteString(s string) error {
	if err := sw.WriteInt32(int32(len(s))); err != nil {
		return err
	}

	_, err := io.WriteString(sw.writer, s)
	return err
}

// WriteDouble encodes a double
func (sw *StreamWriter) WriteDouble(d float64) error {
	value := math.Float64bits(d)
	return sw.WriteInt64(int64(value))
}

// WriteBinary encodes binary
func (sw *StreamWriter) WriteBinary(b []byte) error {
	if err := sw.WriteInt32(int32(len(b))); err != nil {
		return err
	}
	return sw.write(b)
}

// WriteFieldBegin marks the beginning of a new field in a struct. The first
// byte denotes the type and the next two bytes denote the field id.
func (sw *StreamWriter) WriteFieldBegin(f stream.FieldHeader) error {
	// type:1
	if err := sw.writeByte(byte(f.Type)); err != nil {
		return err
	}

	// id:2
	if err := sw.WriteInt16(f.ID); err != nil {
		return err
	}

	return nil
}

// WriteFieldEnd denotes the end of a field. No-op.
func (sw *StreamWriter) WriteFieldEnd() error {
	return nil
}

// WriteStructBegin denotes the beginning of a struct. No-op.
func (sw *StreamWriter) WriteStructBegin() error {
	return nil
}

// WriteStructEnd uses the zero byte to mark the end of a struct.
func (sw *StreamWriter) WriteStructEnd() error {
	return sw.writeByte(0) // end struct
}

// WriteListBegin marks the beginning of a new list. The first byte denotes
// the type of the items and the next four bytes denote the length of the list.
func (sw *StreamWriter) WriteListBegin(l stream.ListHeader) error {
	// vtype:1
	if err := sw.writeByte(byte(l.Type)); err != nil {
		return err
	}

	// length:4
	if err := sw.WriteInt32(int32(l.Length)); err != nil {
		return err
	}

	return nil
}

// WriteListEnd marks the end of a list. No-op.
func (sw *StreamWriter) WriteListEnd() error {
	return nil
}

// WriteSetBegin marks the beginning of a new set. The first byte denotes
// the type of the items and the next four bytes denote the length of the set.
func (sw *StreamWriter) WriteSetBegin(s stream.SetHeader) error {
	// vtype:1
	if err := sw.writeByte(byte(s.Type)); err != nil {
		return err
	}

	// length:4
	if err := sw.WriteInt32(int32(s.Length)); err != nil {
		return err
	}

	return nil
}

// WriteSetEnd marks the end of a set. No-op.
func (sw *StreamWriter) WriteSetEnd() error {
	return nil
}

// WriteMapBegin marks the beginning of a new map. The first byte denotes
// the type of the keys, the second byte denotes the type of the values,
// and the next four bytes denote the length of the map.
func (sw *StreamWriter) WriteMapBegin(m stream.MapHeader) error {
	// ktype:1
	if err := sw.writeByte(byte(m.KeyType)); err != nil {
		return err
	}

	// vtype:1
	if err := sw.writeByte(byte(m.ValueType)); err != nil {
		return err
	}

	// length:4
	if err := sw.WriteInt32(int32(m.Length)); err != nil {
		return err
	}

	return nil
}

// WriteMapEnd marks the end of a map. No-op.
func (sw *StreamWriter) WriteMapEnd() error {
	return nil
}

// Close frees up the resources used by the StreamWriter and returns it back
// to the pool.
func (sw *StreamWriter) Close() error {
	returnStreamWriter(sw)
	return nil
}
