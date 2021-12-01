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

package wire

// ValueList represents a collection of Value objects as an iteration through
// it. This helps us avoid the cost of allocating memory for all collections
// passing through the system.
type ValueList interface {
	// Size returns the size of this lazy list.
	Size() int

	// ValueType specifies the type of values contained in this list.
	ValueType() Type

	// ForEach calls the given function on each element of the list.
	//
	// If any call fails with an error, that error is returned and the
	// iteration is stopped.
	ForEach(f func(Value) error) error

	// Close indicates that the caller is finished reading from the lazy list.
	Close()
}

// MapItemList represents a collection of MapItem objects as an iteration
// through it. This helps us avoid the cost of allocating memory for all
// collections passing through the system.
type MapItemList interface {
	// Size returns the size of this lazy list.
	Size() int

	// KeyType and ValueType specify the kind of values held in this
	// MapItemList.
	KeyType() Type
	ValueType() Type

	// ForEach calls the given function on each element of the list.
	//
	// If any call fails with an error, that error is returned and the
	// iteration is stopped.
	ForEach(f func(MapItem) error) error

	// Close indicates that the caller is finished reading from the lazy list.
	Close()
}

//////////////////////////////////////////////////////////////////////////////

// ValueListFromSlice builds a ValueList from the given slice of Values.
func ValueListFromSlice(t Type, values []Value) ValueList {
	return sliceValueList{t: t, values: values}
}

type sliceValueList struct {
	t      Type
	values []Value
}

func (vs sliceValueList) ValueType() Type {
	return vs.t
}

func (vs sliceValueList) Size() int {
	return len(vs.values)
}

func (vs sliceValueList) ForEach(f func(Value) error) error {
	for _, v := range vs.values {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

func (sliceValueList) Close() {}

//////////////////////////////////////////////////////////////////////////////

// MapItemListFromSlice builds a MapItemList from the given slice of Values.
func MapItemListFromSlice(k, v Type, items []MapItem) MapItemList {
	return sliceMapItemList{ktype: k, vtype: v, items: items}
}

type sliceMapItemList struct {
	ktype, vtype Type
	items        []MapItem
}

func (vs sliceMapItemList) KeyType() Type {
	return vs.ktype
}

func (vs sliceMapItemList) ValueType() Type {
	return vs.vtype
}

func (vs sliceMapItemList) Size() int {
	return len(vs.items)
}

func (vs sliceMapItemList) ForEach(f func(MapItem) error) error {
	for _, v := range vs.items {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

func (sliceMapItemList) Close() {}

//////////////////////////////////////////////////////////////////////////////

// ValueListToSlice builds a slice of values from the given ValueList.
func ValueListToSlice(l ValueList) []Value {
	items := make([]Value, 0, l.Size())
	// explicitly ignoring since we know there will not be an error
	_ = l.ForEach(func(v Value) error {
		items = append(items, v)
		return nil
	})
	return items
}

// MapItemListToSlice builds a slice of values from the given MapItemList.
func MapItemListToSlice(l MapItemList) []MapItem {
	items := make([]MapItem, 0, l.Size())
	// explicitly ignoring since we know there will not be an error
	_ = l.ForEach(func(v MapItem) error {
		items = append(items, v)
		return nil
	})
	return items
}
