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

import "fmt"

// EvaluateValue ensures that the given Value is fully evaluated. Lazy lists
// are spinned and any errors raised by them are returned.
func EvaluateValue(v Value) error {
	switch v.Type() {
	case TBool, TI8, TDouble, TI16, TI32, TI64, TBinary:
		return nil
	case TStruct:
		for _, f := range v.GetStruct().Fields {
			if err := EvaluateValue(f.Value); err != nil {
				return err
			}
		}
		return nil
	case TMap:
		m := v.GetMap()
		defer m.Close()
		return m.ForEach(func(item MapItem) error {
			if err := EvaluateValue(item.Key); err != nil {
				return err
			}
			return EvaluateValue(item.Value)
		})
	case TSet:
		s := v.GetSet()
		defer s.Close()
		return s.ForEach(EvaluateValue)
	case TList:
		l := v.GetList()
		defer l.Close()
		return l.ForEach(EvaluateValue)
	default:
		return fmt.Errorf("unknown type %s", v.Type())
	}
}
