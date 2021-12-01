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
	"io"
	"sync"

	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

var writerPool = sync.Pool{
	New: func() interface{} {
		writer := &Writer{}
		writer.writeValue = writer.WriteValue
		writer.writeMapItem = writer.realWriteMapItem
		return writer
	}}

// Writer implements basic logic for writing the Thrift Binary Protocol to an
// io.Writer.
type Writer struct {
	sw *StreamWriter

	// NOTE:
	// This is a hack to avoid memory allocation in closures. Passing the
	// bound WriteValue or realWriteMapItem methods into a function results in
	// a memory allocation because the system doesn't know we're going to
	// reuse the closure. So we create that bound reference in advance when
	// the writer is created.
	writeValue   func(wire.Value) error
	writeMapItem func(wire.MapItem) error
}

// BorrowWriter fetches a Writer from the system that will write its output to
// the given io.Writer.
//
// This Writer must be returned back using ReturnWriter.
func BorrowWriter(w io.Writer) *Writer {
	streamWriter := NewStreamWriter(w)
	writer := writerPool.Get().(*Writer)
	writer.sw = streamWriter
	return writer
}

// ReturnWriter returns a previously borrowed Writer back to the system.
func ReturnWriter(w *Writer) {
	sw := w.sw
	w.sw = nil
	returnStreamWriter(sw)
	writerPool.Put(w)
}

func (bw *Writer) writeField(f wire.Field) error {
	fh := stream.FieldHeader{
		ID:   f.ID,
		Type: f.Value.Type(),
	}
	if err := bw.sw.WriteFieldBegin(fh); err != nil {
		return err
	}

	// value
	if err := bw.WriteValue(f.Value); err != nil {
		// TODO(abg): Figure out better error handling story. We need access
		// to the underlying error object if it's a network error.
		return fmt.Errorf(
			"failed to write field %d (%v): %s",
			f.ID, f.Value.Type(), err,
		)
	}

	return bw.sw.WriteFieldEnd()
}

func (bw *Writer) writeStruct(s wire.Struct) error {
	if err := bw.sw.WriteStructBegin(); err != nil {
		return err
	}

	for _, f := range s.Fields {
		if err := bw.writeField(f); err != nil {
			return err
		}
	}
	return bw.sw.WriteStructEnd()
}

func (bw *Writer) realWriteMapItem(item wire.MapItem) error {
	if err := bw.WriteValue(item.Key); err != nil {
		return err
	}
	return bw.WriteValue(item.Value)
}

func (bw *Writer) writeMap(m wire.MapItemList) error {
	mh := stream.MapHeader{
		KeyType:   m.KeyType(),
		ValueType: m.ValueType(),
		Length:    m.Size(),
	}
	if err := bw.sw.WriteMapBegin(mh); err != nil {
		return err
	}

	if err := m.ForEach(bw.writeMapItem); err != nil {
		return err
	}

	return bw.sw.WriteMapEnd()
}

func (bw *Writer) writeSet(s wire.ValueList) error {
	sh := stream.SetHeader{
		Type:   s.ValueType(),
		Length: s.Size(),
	}
	if err := bw.sw.WriteSetBegin(sh); err != nil {
		return err
	}

	if err := s.ForEach(bw.writeValue); err != nil {
		return err
	}

	return bw.sw.WriteSetEnd()
}

func (bw *Writer) writeList(l wire.ValueList) error {
	lh := stream.ListHeader{
		Type:   l.ValueType(),
		Length: l.Size(),
	}
	if err := bw.sw.WriteListBegin(lh); err != nil {
		return err
	}

	if err := l.ForEach(bw.writeValue); err != nil {
		return err
	}

	return bw.sw.WriteListEnd()
}

// WriteValue writes the given Thrift value to the underlying stream using the
// Thrift Binary Protocol.
func (bw *Writer) WriteValue(v wire.Value) error {
	switch v.Type() {
	case wire.TBool:
		return bw.sw.WriteBool(v.GetBool())

	case wire.TI8:
		return bw.sw.WriteInt8(v.GetI8())

	case wire.TDouble:
		return bw.sw.WriteDouble(v.GetDouble())

	case wire.TI16:
		return bw.sw.WriteInt16(v.GetI16())

	case wire.TI32:
		return bw.sw.WriteInt32(v.GetI32())

	case wire.TI64:
		return bw.sw.WriteInt64(v.GetI64())

	case wire.TBinary:
		return bw.sw.WriteBinary(v.GetBinary())

	case wire.TStruct:
		return bw.writeStruct(v.GetStruct())

	case wire.TMap:
		return bw.writeMap(v.GetMap())

	case wire.TSet:
		return bw.writeSet(v.GetSet())

	case wire.TList:
		return bw.writeList(v.GetList())

	default:
		return fmt.Errorf("unknown ttype %v", v.Type())
	}
}
