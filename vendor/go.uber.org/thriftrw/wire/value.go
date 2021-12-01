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
	"fmt"
	"math"
	"strings"
)

// An empty []byte with zero length and capacity. We'll use this rather than
// allocating new byte slices for empty []byte.
var _emptyByteSlice = make([]byte, 0)

// Value holds the over-the-wire representation of a Thrift value.
//
// The Type of the value determines which field in the Value is valid.
type Value struct {
	typ Type

	tnumber uint64
	tbinary []byte
	tstruct Struct
	tcoll   interface{} // set/map/list
}

// Type retrieves the type of value inside a Value.
func (v *Value) Type() Type {
	return v.typ
}

// Get retrieves whatever value the given Value contains.
func (v *Value) Get() interface{} {
	switch v.typ {
	case TBool:
		return v.GetBool()
	case TI8:
		return v.GetI8()
	case TDouble:
		return v.GetDouble()
	case TI16:
		return v.GetI16()
	case TI32:
		return v.GetI32()
	case TI64:
		return v.GetI64()
	case TBinary:
		return v.GetBinary()
	case TStruct:
		return v.GetStruct()
	case TMap:
		return v.GetMap()
	case TSet, TList:
		return v.GetList()
	default:
		panic(fmt.Sprintf("Unknown value type %v", v.typ))
	}
}

// NewValueBool constructs a new Value that contains a boolean.
func NewValueBool(v bool) Value {
	n := uint64(0)
	if v {
		n = 1
	}
	return Value{
		typ:     TBool,
		tnumber: n,
	}
}

// GetBool gets the Bool value from a Value.
func (v *Value) GetBool() bool {
	return v.tnumber != 0
}

// NewValueI8 constructs a new Value that contains a byte
func NewValueI8(v int8) Value {
	return Value{
		typ:     TI8,
		tnumber: uint64(v),
	}
}

// GetI8 gets the I8 value from a Value.
func (v *Value) GetI8() int8 {
	return int8(v.tnumber)
}

// NewValueDouble constructs a new Value that contains a double.
func NewValueDouble(v float64) Value {
	return Value{
		typ:     TDouble,
		tnumber: math.Float64bits(v),
	}
}

// GetDouble gets the Double value from a Value.
func (v *Value) GetDouble() float64 {
	return math.Float64frombits(v.tnumber)
}

// NewValueI16 constructs a new Value that contains a 16-bit integer.
func NewValueI16(v int16) Value {
	return Value{
		typ:     TI16,
		tnumber: uint64(v),
	}
}

// GetI16 gets the I16 value from a Value.
func (v *Value) GetI16() int16 {
	return int16(v.tnumber)
}

// NewValueI32 constructs a new Value that contains a 32-bit integer.
func NewValueI32(v int32) Value {
	return Value{
		typ:     TI32,
		tnumber: uint64(v),
	}
}

// GetI32 gets the I32 value from a Value.
func (v *Value) GetI32() int32 {
	return int32(v.tnumber)
}

// NewValueI64 constructs a new Value that contains a 64-bit integer.
func NewValueI64(v int64) Value {
	return Value{
		typ:     TI64,
		tnumber: uint64(v),
	}
}

// GetI64 gets the I64 value from a Value.
func (v *Value) GetI64() int64 {
	return int64(v.tnumber)
}

// NewValueBinary constructs a new Value that contains a binary string.
func NewValueBinary(v []byte) Value {
	if v == nil {
		v = _emptyByteSlice
	}
	return Value{
		typ:     TBinary,
		tbinary: v,
	}
}

// NewValueString constructs a new Value that contains a string.
func NewValueString(v string) Value {
	return Value{
		typ:     TBinary,
		tbinary: unsafeStringToBytes(v),
	}
}

// GetBinary gets the Binary value from a Value.
func (v *Value) GetBinary() []byte {
	return v.tbinary
}

// GetString gets a string value from a Value.
func (v *Value) GetString() string {
	return unsafeBytesToString(v.tbinary)
}

// NewValueStruct constructs a new Value that contains a struct.
func NewValueStruct(v Struct) Value {
	return Value{
		typ:     TStruct,
		tstruct: v,
	}
}

// GetStruct gets the Struct value from a Value.
func (v *Value) GetStruct() Struct {
	return v.tstruct
}

// NewValueMap constructs a new Value that contains a map.
func NewValueMap(v MapItemList) Value {
	return Value{
		typ:   TMap,
		tcoll: v,
	}
}

// GetMap gets the Map value from a Value.
func (v *Value) GetMap() MapItemList {
	return v.tcoll.(MapItemList)
}

// NewValueSet constructs a new Value that contains a set.
func NewValueSet(v ValueList) Value {
	return Value{
		typ:   TSet,
		tcoll: v,
	}
}

// GetSet gets the Set value from a Value.
func (v *Value) GetSet() ValueList {
	return v.tcoll.(ValueList)
}

// NewValueList constructs a new Value that contains a list.
func NewValueList(v ValueList) Value {
	return Value{
		typ:   TList,
		tcoll: v,
	}
}

// GetList gets the List value from a Value.
func (v *Value) GetList() ValueList {
	return v.tcoll.(ValueList)
}

func (v Value) String() string {
	switch v.typ {
	case TBool:
		return fmt.Sprintf("TBool(%v)", v.GetBool())
	case TI8:
		return fmt.Sprintf("TI8(%v)", v.GetI8())
	case TDouble:
		return fmt.Sprintf("TDouble(%v)", v.GetDouble())
	case TI16:
		return fmt.Sprintf("TI16(%v)", v.GetI16())
	case TI32:
		return fmt.Sprintf("TI32(%v)", v.GetI32())
	case TI64:
		return fmt.Sprintf("TI64(%v)", v.GetI64())
	case TBinary:
		return fmt.Sprintf("TBinary(%v)", v.tbinary)
	case TStruct:
		return fmt.Sprintf("TStruct(%v)", v.tstruct)
	case TMap:
		return fmt.Sprintf("TMap(%v)", v.tcoll)
	case TSet:
		return fmt.Sprintf("TSet(%v)", v.tcoll)
	case TList:
		return fmt.Sprintf("TList(%v)", v.tcoll)
	default:
		panic(fmt.Sprintf("Unknown value type %v", v.typ))
	}
}

// Struct provides a wire-level representation of a struct.
//
// At this level, structs don't have names or named fields.
type Struct struct {
	Fields []Field
}

func (s Struct) fieldMap() map[int16]Value {
	m := make(map[int16]Value, len(s.Fields))
	for _, f := range s.Fields {
		m[f.ID] = f.Value
	}
	return m
}

func (s Struct) String() string {
	fields := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		fields[i] = field.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(fields, ", "))
}

// Field is a single field inside a Struct.
type Field struct {
	ID    int16
	Value Value
}

func (f Field) String() string {
	return fmt.Sprintf("%v: %v", f.ID, f.Value)
}

// MapItem is a single item in a Map.
type MapItem struct {
	Key   Value
	Value Value
}

func (mi MapItem) String() string {
	return fmt.Sprintf("%v: %v", mi.Key, mi.Value)
}
