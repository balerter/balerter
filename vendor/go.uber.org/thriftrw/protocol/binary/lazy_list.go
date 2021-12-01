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
	"sync"

	"go.uber.org/thriftrw/wire"
)

var (
	lazyValueListPool = sync.Pool{New: func() interface{} {
		return &lazyValueList{}
	}}
	lazyMapItemListPool = sync.Pool{New: func() interface{} {
		return &lazyMapItemList{}
	}}
)

func borrowLazyValueList() *lazyValueList {
	return lazyValueListPool.Get().(*lazyValueList)
}

func borrowLazyMapItemList() *lazyMapItemList {
	return lazyMapItemListPool.Get().(*lazyMapItemList)
}

// lazyValueList is an implementation of ValueList which parses Values from a
// Reader on-demand.
type lazyValueList struct {
	count       int32
	typ         wire.Type
	readerAt    io.ReaderAt
	startOffset int64
}

func (ll *lazyValueList) ValueType() wire.Type {
	return ll.typ
}

func (ll *lazyValueList) Size() int {
	return int(ll.count)
}

func (ll *lazyValueList) ForEach(f func(wire.Value) error) error {
	off := ll.startOffset
	reader := newReader(ll.readerAt, off)
	defer reader.close()

	for i := int32(0); i < ll.count; i++ {
		var (
			val wire.Value
			err error
		)

		val, off, err = reader.ReadValue(ll.typ, off)
		if err != nil {
			return err
		}

		if err := f(val); err != nil {
			return err
		}
	}
	return nil
}

func (ll *lazyValueList) Close() {
	ll.readerAt = nil
	lazyValueListPool.Put(ll)
}

// lazyMapItemList is an implementation of MapItemList which parses MapItems
// from a Reader on-demand.
type lazyMapItemList struct {
	ktype, vtype wire.Type
	count        int32
	readerAt     io.ReaderAt
	startOffset  int64
}

func (lm *lazyMapItemList) KeyType() wire.Type {
	return lm.ktype
}

func (lm *lazyMapItemList) ValueType() wire.Type {
	return lm.vtype
}

func (lm *lazyMapItemList) Size() int {
	return int(lm.count)
}

func (lm *lazyMapItemList) ForEach(f func(wire.MapItem) error) error {
	off := lm.startOffset
	reader := newReader(lm.readerAt, off)
	defer reader.close()

	for i := int32(0); i < lm.count; i++ {
		var (
			k, v wire.Value
			err  error
		)

		k, off, err = reader.ReadValue(lm.ktype, off)
		if err != nil {
			return err
		}

		v, off, err = reader.ReadValue(lm.vtype, off)
		if err != nil {
			return err
		}

		item := wire.MapItem{Key: k, Value: v}
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

func (lm *lazyMapItemList) Close() {
	lm.readerAt = nil
	lazyMapItemListPool.Put(lm)
}
