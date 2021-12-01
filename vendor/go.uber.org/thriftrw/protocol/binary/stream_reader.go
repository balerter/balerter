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
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"sync"

	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

// Requests for byte slices longer than this will use a dynamically resizing
// buffer.
const bytesAllocThreshold = 1048576 // 1 MB

// For the reader, we keep track of the read offset manually everywhere so
// that we can implement lazy collections without extra allocations

// fixedWidth returns the encoded size of a value of the given type. If the
// type's width depends on the value, -1 is returned.
func fixedWidth(t wire.Type) int64 {
	switch t {
	case wire.TBool:
		return 1
	case wire.TI8:
		return 1
	case wire.TDouble:
		return 8
	case wire.TI16:
		return 2
	case wire.TI32:
		return 4
	case wire.TI64:
		return 8
	default:
		return -1
	}
}

// StreamReader provides an implementation of a "stream.Reader".
type StreamReader struct {
	reader io.Reader
	buffer [8]byte

	// discard points to either discardOffset or discardStream based on
	// the implementation of the io.Reader we're using.
	discard func(int64) error

	// These are bound versions of the discardStream and discardSeek
	// methods on the StreamReader. Putting them here ensures that we don't
	// cause an alloc when we do "sr.discard = sr.discardOffset".
	_discardStream func(int64) error
	_discardSeek   func(int64) error

	// This field is set only if the wrapped reader is an io.Seeker. ONLY
	// USE if you are discardSeek.
	_seeker io.Seeker
}

var streamReaderPool = sync.Pool{
	New: func() interface{} {
		sr := new(StreamReader)
		sr._discardSeek = sr.discardSeek
		sr._discardStream = sr.discardStream
		return sr
	},
}

// NewStreamReader fetches a StreamReader from the system that will write
// its output to the given io.Reader.
//
// This StreamReader must be closed using `Close()`
func NewStreamReader(r io.Reader) *StreamReader {
	sr := streamReaderPool.Get().(*StreamReader)
	sr.reader = r
	sr.discard = sr._discardStream
	if seeker, ok := r.(io.Seeker); ok {
		// If we're wrapping a seeker (like *offsetReader), we can skip
		// bytes much more efficiently.
		sr._seeker = seeker
		sr.discard = sr._discardSeek
	}
	return sr
}

func returnStreamReader(sr *StreamReader) {
	sr.reader = nil
	sr._seeker = nil
	streamReaderPool.Put(sr)
}

func (sr *StreamReader) read(bs []byte) (int, error) {
	n, err := io.ReadFull(sr.reader, bs)

	if err == io.EOF {
		// All EOFs are unexpected when streaming
		err = io.ErrUnexpectedEOF
	}

	return n, err
}

func (sr *StreamReader) discardSeek(n int64) error {
	_, err := sr._seeker.Seek(n, io.SeekCurrent)
	return err
}

func (sr *StreamReader) discardStream(n int64) error {
	_, err := io.CopyN(ioutil.Discard, sr.reader, n)
	if err == io.EOF {
		// All EOFs are unexpected when streaming
		err = io.ErrUnexpectedEOF
	}

	return err
}

// ReadBool reads a Thrift encoded bool value, returning a bool.
func (sr *StreamReader) ReadBool() (bool, error) {
	bs := sr.buffer[0:1]
	if _, err := sr.read(bs); err != nil {
		return false, err
	}

	switch bs[0] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, decodeErrorf("invalid bool value: %v", bs[0])
	}
}

// ReadInt8 reads a Thrift encoded int8 value.
func (sr *StreamReader) ReadInt8() (int8, error) {
	bs := sr.buffer[0:1]
	_, err := sr.read(bs)
	return int8(bs[0]), err
}

// ReadInt16 reads a Thrift encoded int16 value.
func (sr *StreamReader) ReadInt16() (int16, error) {
	bs := sr.buffer[0:2]
	_, err := sr.read(bs)
	return int16(bigEndian.Uint16(bs)), err
}

// ReadInt32 reads a Thrift encoded int32 value.
func (sr *StreamReader) ReadInt32() (int32, error) {
	bs := sr.buffer[0:4]
	_, err := sr.read(bs)
	return int32(bigEndian.Uint32(bs)), err
}

// ReadInt64 reads a Thrift encoded int64 value.
func (sr *StreamReader) ReadInt64() (int64, error) {
	bs := sr.buffer[0:8]
	_, err := sr.read(bs)
	return int64(bigEndian.Uint64(bs)), err
}

// ReadString reads a Thrift encoded string.
func (sr *StreamReader) ReadString() (string, error) {
	bs, err := sr.ReadBinary()
	return string(bs), err
}

// ReadDouble reads a Thrift encoded double, returning a float64.
func (sr *StreamReader) ReadDouble() (float64, error) {
	val, err := sr.ReadInt64()
	return math.Float64frombits(uint64(val)), err
}

// ReadBinary reads a Thrift encoded binary type, returning a byte array.
func (sr *StreamReader) ReadBinary() ([]byte, error) {
	length, err := sr.ReadInt32()
	if err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, decodeErrorf("negative length %v specified for binary field", length)
	}

	if length == 0 {
		return []byte{}, nil
	}

	if length > bytesAllocThreshold {
		var buf bytes.Buffer
		_, err := io.CopyN(&buf, sr.reader, int64(length))
		if err == io.EOF {
			// All EOFs are unexpected when streaming
			err = io.ErrUnexpectedEOF
		}

		return buf.Bytes(), err
	}

	bs := make([]byte, length)
	_, err = sr.read(bs)
	return bs, err
}

// ReadStructBegin reads the "beginning" of a Thrift encoded struct.  Since
// there is no encoding for the beginning of a struct, this is a noop.
func (sr *StreamReader) ReadStructBegin() error {
	return nil
}

// ReadStructEnd reads the "end" of a Thrift encoded struct.  Since
// `ReadFieldBegin` will already be interpreting field-type of whether it's a
// stop field or not, there is no real representation of a struct's end, making
// this a noop.
func (sr *StreamReader) ReadStructEnd() error {
	return nil
}

// ReadFieldBegin reads off a Thrift encoded field-header returning that and a
// 'bool' representing whether or not a field-value follows.
// A 'false' without any error means that it has reached the stop-field.  There
// is no guarantee that the field-header is valid in this case.
func (sr *StreamReader) ReadFieldBegin() (fh stream.FieldHeader, ok bool, err error) {
	fieldType, err := sr.ReadInt8()
	if err != nil {
		return fh, false, err
	}

	// fieldType 0 is the equivalent of the stop field, signaling the end
	if fieldType == 0 {
		return fh, false, nil
	}

	fieldID, err := sr.ReadInt16()
	if err != nil {
		return fh, false, err
	}

	fh.ID = fieldID
	fh.Type = wire.Type(fieldType)
	return fh, true, nil
}

// ReadFieldEnd reads the "end" of a Thrift encoded field  Since there is no
// encoding for the end of a field, this is a noop.
func (sr *StreamReader) ReadFieldEnd() error {
	return nil
}

// ReadListBegin reads off the list header of a Thrift encoded list.
func (sr *StreamReader) ReadListBegin() (lh stream.ListHeader, err error) {
	elemType, listSize, err := sr.readTypeSizeHeader()
	if err != nil {
		return lh, err
	}

	lh.Type = wire.Type(elemType)
	lh.Length = int(listSize)
	return lh, nil
}

// ReadListEnd reads the "end" of a Thrift encoded list.  Since there is no
// encoding for the end of a list, this is a noop.
func (sr *StreamReader) ReadListEnd() error {
	return nil
}

// ReadSetBegin reads off the set header of a Thrift encoded set.
func (sr *StreamReader) ReadSetBegin() (sh stream.SetHeader, err error) {
	elemType, setSize, err := sr.readTypeSizeHeader()
	if err != nil {
		return sh, err
	}

	sh.Type = elemType
	sh.Length = setSize
	return sh, nil
}

// ReadSetEnd reads the "end" of a Thrift encoded list.  Since there is no
// encoding for the end of a set, this is a noop.
func (sr *StreamReader) ReadSetEnd() error {
	return nil
}

func (sr *StreamReader) readTypeSizeHeader() (wire.Type, int, error) {
	elemType, err := sr.ReadInt8()
	if err != nil {
		return 0, 0, err
	}

	size, err := sr.ReadInt32()
	if err != nil {
		return 0, 0, err
	}

	if size < 0 {
		return 0, 0, decodeErrorf("got negative length: %v", size)
	}

	return wire.Type(elemType), int(size), nil
}

// ReadMapBegin reads off the map header of a Thrift encoded map.
func (sr *StreamReader) ReadMapBegin() (mh stream.MapHeader, err error) {
	keyType, err := sr.ReadInt8()
	if err != nil {
		return mh, err
	}

	valueType, err := sr.ReadInt8()
	if err != nil {
		return mh, err
	}

	size, err := sr.ReadInt32()
	if err != nil {
		return mh, err
	}

	if size < 0 {
		return mh, decodeErrorf("got negative length: %v", size)
	}

	mh.KeyType = wire.Type(keyType)
	mh.ValueType = wire.Type(valueType)
	mh.Length = int(size)
	return mh, nil
}

// ReadMapEnd reads the "end" of a Thrift encoded list.  Since there is no
// encoding for the end of a map, this is a noop.
func (sr *StreamReader) ReadMapEnd() error {
	return nil
}

// Skip skips fully over the provided Thrift type.
func (sr *StreamReader) Skip(t wire.Type) error {
	// if it's a primitive type, skip over it based on its width
	if w := fixedWidth(t); w > 0 {
		return sr.discard(w)
	}

	switch t {
	case wire.TBinary:
		length, err := sr.ReadInt32()
		if err != nil {
			return err
		}

		if length < 0 {
			return decodeErrorf("got negative length: %v", length)
		}

		return sr.discard(int64(length))
	case wire.TStruct:
		return sr.skipStruct()
	case wire.TMap:
		return sr.skipMap()
	case wire.TSet:
		return sr.skipList()
	case wire.TList:
		return sr.skipList()
	default:
		return decodeErrorf("unknown ttype %v", t)
	}
}

// Close frees up the resources used by the StreamReader and returns it back
// to the pool.
func (sr *StreamReader) Close() error {
	returnStreamReader(sr)
	return nil
}

func (sr *StreamReader) skipStruct() error {
	fieldType, err := sr.ReadInt8()
	if err != nil {
		return err
	}

	for fieldType != 0 {
		// field id
		if err := sr.discard(int64(2)); err != nil {
			return err
		}

		if err := sr.Skip(wire.Type(fieldType)); err != nil {
			return err
		}

		if fieldType, err = sr.ReadInt8(); err != nil {
			return err
		}
	}

	return nil
}

func (sr *StreamReader) skipMap() error {
	key, err := sr.ReadInt8()
	if err != nil {
		return err
	}

	value, err := sr.ReadInt8()
	if err != nil {
		return err
	}

	size, err := sr.ReadInt32()
	if err != nil {
		return err
	}

	if size < 0 {
		return decodeErrorf("got negative length: %v", size)
	}

	return sr.skipMapItems(wire.Type(key), wire.Type(value), size)
}

func (sr *StreamReader) skipMapItems(key, value wire.Type, size int32) error {
	keyWidth := fixedWidth(key)
	valueWidth := fixedWidth(value)
	if keyWidth > 0 && valueWidth > 0 {
		length := int64(size) * (keyWidth + valueWidth)
		return sr.discard(length)
	}

	for i := int32(0); i < size; i++ {
		if err := sr.Skip(key); err != nil {
			return err
		}

		if err := sr.Skip(value); err != nil {
			return err
		}
	}

	return nil
}

func (sr *StreamReader) skipList() error {
	elemType, size, err := sr.readTypeSizeHeader()
	if err != nil {
		return err
	}

	return sr.skipListItems(elemType, size)
}

func (sr *StreamReader) skipListItems(elemType wire.Type, size int) error {
	width := fixedWidth(elemType)
	if width > 0 {
		length := width * int64(size)
		return sr.discard(length)
	}

	for i := 0; i < size; i++ {
		if err := sr.Skip(elemType); err != nil {
			return err
		}
	}

	return nil
}
