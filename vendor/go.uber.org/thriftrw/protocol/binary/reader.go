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

	"go.uber.org/thriftrw/wire"
)

// offsetReader provides a type that satisfies an io.Reader with only an
// io.ReaderAt.
type offsetReader struct {
	offset int64
	reader io.ReaderAt
}

var (
	_ io.Reader = (*offsetReader)(nil)
	_ io.Seeker = (*offsetReader)(nil)
)

// Read reads len(p) bytes into p.
func (or *offsetReader) Read(p []byte) (int, error) {
	n, err := or.reader.ReadAt(p, or.offset)
	or.offset += int64(n)

	return n, err
}

func (or *offsetReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		or.offset = offset
	case io.SeekCurrent:
		or.offset += offset
	default:
		return or.offset, fmt.Errorf("unsupported whence %d", whence)
	}
	return or.offset, nil
}

// reader functions as the actual reader behind the exported `Reader` type.
// This is necessary to avoid new calls to a `Reader.ReadValue` from changing
// the offset in already running 'ReadValue' calls.
type reader struct {
	or *offsetReader
	sr *StreamReader
}

func newReader(r io.ReaderAt, off int64) reader {
	or := offsetReader{reader: r, offset: off}

	return reader{
		or: &or,
		sr: NewStreamReader(&or),
	}
}

func (r *reader) readStructStream() (wire.Struct, error) {
	var fields []wire.Field

	if err := r.sr.ReadStructBegin(); err != nil {
		return wire.Struct{}, err
	}

	fh, ok, err := r.sr.ReadFieldBegin()
	if err != nil {
		return wire.Struct{}, err
	}

	for ok {
		val, _, err := r.ReadValue(fh.Type, r.or.offset)
		if err != nil {
			return wire.Struct{}, err
		}

		fields = append(fields, wire.Field{ID: fh.ID, Value: val})
		if err := r.sr.ReadFieldEnd(); err != nil {
			return wire.Struct{}, err
		}

		if fh, ok, err = r.sr.ReadFieldBegin(); err != nil {
			return wire.Struct{}, err
		}
	}

	if err := r.sr.ReadStructEnd(); err != nil {
		return wire.Struct{}, err
	}

	return wire.Struct{Fields: fields}, nil
}

func (r *reader) readMapStream() (wire.MapItemList, error) {
	mh, err := r.sr.ReadMapBegin()
	if err != nil {
		return nil, err
	}

	start := r.or.offset
	if err := r.sr.skipMapItems(mh.KeyType, mh.ValueType, int32(mh.Length)); err != nil {
		return nil, err
	}

	if err := r.sr.ReadMapEnd(); err != nil {
		return nil, err
	}

	items := borrowLazyMapItemList()
	items.ktype = mh.KeyType
	items.vtype = mh.ValueType
	items.count = int32(mh.Length)
	items.readerAt = r.or.reader
	items.startOffset = start

	return items, nil
}

func (r *reader) readListStream() (wire.ValueList, error) {
	lh, err := r.sr.ReadListBegin()
	if err != nil {
		return nil, err
	}

	start := r.or.offset
	if err := r.sr.skipListItems(lh.Type, lh.Length); err != nil {
		return nil, err
	}

	if err := r.sr.ReadListEnd(); err != nil {
		return nil, err
	}

	items := borrowLazyValueList()
	items.count = int32(lh.Length)
	items.typ = lh.Type
	items.readerAt = r.or.reader
	items.startOffset = start

	return items, nil
}

func (r *reader) readSetStream() (wire.ValueList, error) {
	sh, err := r.sr.ReadSetBegin()
	if err != nil {
		return nil, err
	}

	start := r.or.offset
	if err := r.sr.skipListItems(sh.Type, sh.Length); err != nil {
		return nil, err
	}

	if err := r.sr.ReadSetEnd(); err != nil {
		return nil, err
	}

	items := borrowLazyValueList()
	items.count = int32(sh.Length)
	items.typ = sh.Type
	items.readerAt = r.or.reader
	items.startOffset = start

	return items, nil
}

func (r *reader) close() error {
	err := r.sr.Close()
	r.sr = nil
	r.or = nil

	return err
}

// ReadValue is the underlying call made from the exported `Reader.ReadValue`
// that's meant to be safe for concurrent calls.
func (r *reader) ReadValue(t wire.Type, off int64) (wire.Value, int64, error) {
	r.or.offset = off

	switch t {
	case wire.TBool:
		b, err := r.sr.ReadBool()
		return wire.NewValueBool(b), r.or.offset, err

	case wire.TI8:
		b, err := r.sr.ReadInt8()
		return wire.NewValueI8(int8(b)), r.or.offset, err

	case wire.TDouble:
		value, err := r.sr.ReadDouble()
		return wire.NewValueDouble(value), r.or.offset, err

	case wire.TI16:
		n, err := r.sr.ReadInt16()
		return wire.NewValueI16(n), r.or.offset, err

	case wire.TI32:
		n, err := r.sr.ReadInt32()
		return wire.NewValueI32(n), r.or.offset, err

	case wire.TI64:
		n, err := r.sr.ReadInt64()
		return wire.NewValueI64(n), r.or.offset, err

	case wire.TBinary:
		v, err := r.sr.ReadBinary()
		return wire.NewValueBinary(v), r.or.offset, err

	case wire.TStruct:
		s, err := r.readStructStream()
		return wire.NewValueStruct(s), r.or.offset, err

	case wire.TMap:
		m, err := r.readMapStream()
		return wire.NewValueMap(m), r.or.offset, err

	case wire.TSet:
		s, err := r.readSetStream()
		return wire.NewValueSet(s), r.or.offset, err

	case wire.TList:
		l, err := r.readListStream()
		return wire.NewValueList(l), r.or.offset, err

	default:
		return wire.Value{}, r.or.offset, decodeErrorf("unknown ttype %v", t)
	}
}

// Reader implements a parser for the Thrift Binary Protocol based on an
// io.ReaderAt.
type Reader struct {
	reader io.ReaderAt
}

// NewReader builds a new Reader based on the given io.ReaderAt.
func NewReader(r io.ReaderAt) Reader {
	return Reader{reader: r}
}

// ReadValue reads a value off the given type off the wire starting at the
// given offset.
//
// Returns the Value, the new offset, and an error if there was a decode error.
func (br *Reader) ReadValue(t wire.Type, off int64) (wire.Value, int64, error) {
	reader := newReader(br.reader, off)
	defer reader.close()
	return reader.ReadValue(t, off)
}
