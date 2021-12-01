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

import (
	"bytes"
	"errors"
	"fmt"
)

// errNotEquals is a sentinel error used while iterating through ValueLists to
// indicate that two values did not match.
var errNotEquals = errors.New("values are not equal")

// ValuesAreEqual checks if two values are equal.
func ValuesAreEqual(left, right Value) bool {
	if left.typ != right.typ {
		return false
	}

	switch left.typ {
	case TBool:
		return left.GetBool() == right.GetBool()
	case TI8:
		return left.GetI8() == right.GetI8()
	case TDouble:
		return left.GetDouble() == right.GetDouble()
	case TI16:
		return left.GetI16() == right.GetI16()
	case TI32:
		return left.GetI32() == right.GetI32()
	case TI64:
		return left.GetI64() == right.GetI64()
	case TBinary:
		return bytes.Equal(left.tbinary, right.tbinary)
	case TStruct:
		return StructsAreEqual(left.tstruct, right.tstruct)
	case TMap:
		return MapsAreEqual(left.tcoll.(MapItemList), right.tcoll.(MapItemList))
	case TSet:
		return SetsAreEqual(left.tcoll.(ValueList), right.tcoll.(ValueList))
	case TList:
		return ListsAreEqual(left.tcoll.(ValueList), right.tcoll.(ValueList))
	default:
		return false
	}
}

// StructsAreEqual checks if two structs are equal.
func StructsAreEqual(left, right Struct) bool {
	if len(left.Fields) != len(right.Fields) {
		return false
	}

	// Fields are unordered so we need to build a map to actually compare
	// them.

	leftFields := left.fieldMap()
	rightFields := right.fieldMap()

	for i, lvalue := range leftFields {
		if rvalue, ok := rightFields[i]; !ok {
			return false
		} else if !ValuesAreEqual(lvalue, rvalue) {
			return false
		}
	}

	return true
}

// SetsAreEqual checks if two sets are equal.
func SetsAreEqual(left, right ValueList) bool {
	if left.ValueType() != right.ValueType() {
		return false
	}
	if left.Size() != right.Size() {
		return false
	}

	if isHashable(left.ValueType()) {
		return setsArEqualHashable(left.Size(), left, right)
	}
	return setsAreEqualUnhashable(left.Size(), left, right)
}

// setsArEqualHashable checks if two unordered ValueLists are equal, provided
// that they contain items that are hashable -- that is, the items can be used
// as keys in a map.
func setsArEqualHashable(size int, l, r ValueList) bool {
	m := make(map[interface{}]bool, size)
	// explicitly ignoring since we know there will not be an error
	_ = l.ForEach(func(v Value) error {
		m[toHashable(v)] = true
		return nil
	})

	return errNotEquals != r.ForEach(func(v Value) error {
		if _, ok := m[toHashable(v)]; !ok {
			return errNotEquals
		}
		return nil
	})
}

// setsAreEqualUnhashable checks if two unordered ValueLists are equal for
// types that are not hashable. Note that this is O(n^2) in time complexity.
func setsAreEqualUnhashable(size int, l, r ValueList) bool {
	lItems := ValueListToSlice(l)

	return errNotEquals != r.ForEach(func(rItem Value) error {
		matched := false
		for _, lItem := range lItems {
			if ValuesAreEqual(lItem, rItem) {
				matched = true
				break
			}
		}
		if !matched {
			return errNotEquals
		}
		return nil
	})
}

// MapsAreEqual checks if two maps are equal.
func MapsAreEqual(left, right MapItemList) bool {
	if left.KeyType() != right.KeyType() {
		return false
	}
	if left.ValueType() != right.ValueType() {
		return false
	}
	if left.Size() != right.Size() {
		return false
	}

	if isHashable(left.KeyType()) {
		return mapsAreEqualHashable(left.Size(), left, right)
	}
	return mapsAreEqualUnhashable(left.Size(), left, right)
}

func mapsAreEqualHashable(size int, l, r MapItemList) bool {
	m := make(map[interface{}]Value, size)

	// explicitly ignoring since we know there will not be an error
	_ = l.ForEach(func(item MapItem) error {
		m[toHashable(item.Key)] = item.Value
		return nil
	})

	return errNotEquals != r.ForEach(func(item MapItem) error {
		lValue, ok := m[toHashable(item.Key)]
		if !ok {
			return errNotEquals
		}
		if !ValuesAreEqual(lValue, item.Value) {
			return errNotEquals
		}
		return nil
	})
}

func mapsAreEqualUnhashable(size int, l, r MapItemList) bool {
	lItems := MapItemListToSlice(l)

	return errNotEquals != r.ForEach(func(rItem MapItem) error {
		matched := false
		for _, lItem := range lItems {
			if !ValuesAreEqual(lItem.Key, rItem.Key) {
				continue
			}
			if !ValuesAreEqual(lItem.Value, rItem.Value) {
				continue
			}
			matched = true
		}

		if !matched {
			return errNotEquals
		}
		return nil
	})
}

func isHashable(t Type) bool {
	switch t {
	case TBool, TI8, TDouble, TI16, TI32, TI64, TBinary:
		return true
	default:
		return false
	}
}

func toHashable(v Value) interface{} {
	switch v.Type() {
	case TBool, TI8, TDouble, TI16, TI32, TI64:
		return v.Get()
	case TBinary:
		return string(v.GetBinary())
	default:
		panic(fmt.Sprintf("value is not hashable: %v", v))
	}
}

// ListsAreEqual checks if two lists are equal.
func ListsAreEqual(left, right ValueList) bool {
	if left.ValueType() != right.ValueType() {
		return false
	}
	if left.Size() != right.Size() {
		return false
	}

	leftItems := ValueListToSlice(left)
	rightItems := ValueListToSlice(right)

	for i, lv := range leftItems {
		rv := rightItems[i]
		if !ValuesAreEqual(lv, rv) {
			return false
		}
	}

	return true
}
