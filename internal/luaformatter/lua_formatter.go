package luaformatter

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
	"sort"
)

var (
	ErrUnsupportedType         = errors.New("unsupported type")
	ErrTableIsNil              = errors.New("table is nil")
	ErrKeyMustBeAString        = errors.New("key must be a string")
	ErrWrongValue              = errors.New("value must be a string, number, boolean or table")
	ErrUnexpectedDataLoss      = errors.New("unexpected data loss")
	ErrUnexpectedDuplicatedKey = errors.New("unexpected duplicated key")
)

const (
	errorMarshalingString = "__ERROR_MARSHALING__"
)

func ValuesToStringNoErr(values []lua.LValue) string {
	v, err := ValuesToString(values)
	if err != nil {
		return errorMarshalingString
	}

	return v
}

func ValueToStringNoErr(value lua.LValue) string {
	v, err := ValueToString(value)
	if err != nil {
		return errorMarshalingString
	}

	return v
}

func ValuesToString(values []lua.LValue) (string, error) {
	s := "["
	for _, v := range values {
		ss, err := ValueToString(v)
		if err != nil {
			return "", err
		}
		s += ss + ","
	}
	s = s[:len(s)-1] + "]"

	return s, nil
}

func ValueToString(v lua.LValue) (string, error) {
	switch v.Type() {
	case lua.LTTable:
		return TableToString(v.(*lua.LTable))
	case lua.LTString, lua.LTNumber, lua.LTBool:
		return v.String(), nil
	}

	return "", ErrUnsupportedType
}

func TableToString(table *lua.LTable) (string, error) {
	if table == nil {
		return "", ErrTableIsNil
	}

	return toString(table)
}

func toString(table *lua.LTable) (string, error) {
	var err error

	data := make(map[string]string)
	var keys []string

	table.ForEach(func(value lua.LValue, value2 lua.LValue) {
		if value.Type() != lua.LTString {
			err = ErrKeyMustBeAString
			return
		}

		key := value.String()
		if _, ok := data[key]; ok {
			err = ErrUnexpectedDuplicatedKey
			return
		}
		keys = append(keys, key)
		var v string

		switch value2.Type() {
		case lua.LTString:
			v = "\"" + value2.String() + "\""
		case lua.LTNumber, lua.LTBool:
			v = value2.String()
		case lua.LTTable:
			v, err = toString(value2.(*lua.LTable))
			if err != nil {
				return
			}
		default:
			err = ErrWrongValue
			return
		}

		data[key] = v
	})

	if err != nil {
		return "", err
	}

	// Sorting keys for same order for any cases
	sort.Strings(keys)

	if len(keys) == 0 {
		return "{}", nil
	}

	s := "{"
	for _, key := range keys {
		s += "\"" + key + "\":"
		v, ok := data[key]
		if !ok {
			err = ErrUnexpectedDataLoss
			return "", err
		}
		s += v + ","
	}

	s = s[:len(s)-1] + "}"

	return s, err
}
