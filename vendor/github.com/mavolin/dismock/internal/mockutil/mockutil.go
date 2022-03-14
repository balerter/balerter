package mockutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WriteJSON writes the passed value to the passed http.ResponseWriter.
func WriteJSON(t *testing.T, w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	require.NoError(t, err)
}

var (
	nullableBool   = reflect.TypeOf(new(option.NullableBoolData))
	nullableUint   = reflect.TypeOf(new(option.NullableUintData))
	nullableInt    = reflect.TypeOf(new(option.NullableIntData))
	nullableString = reflect.TypeOf(new(option.NullableStringData))
	nullableColor  = reflect.TypeOf(new(option.NullableColorData))

	nilBool   = reflect.ValueOf((*option.NullableBoolData)(nil))
	nilUint   = reflect.ValueOf((*option.NullableUintData)(nil))
	nilInt    = reflect.ValueOf((*option.NullableIntData)(nil))
	nilString = reflect.ValueOf((*option.NullableStringData)(nil))
	nilColor  = reflect.ValueOf((*option.NullableColorData)(nil))
)

// CheckJSON checks the body extracted from the passed io.ReadCloser against
// the passed expected value, assuming the body contains JSON data.
// v will be used to decode into and should not contain any data.
func CheckJSON(t *testing.T, r io.ReadCloser, v interface{}, expect interface{}) {
	checkJSON(t, r, v, expect)

	require.NoError(t, r.Close())
}

// CheckMultipart checks if the passed request contains the passed JSON body
// and the passed []api.SendMessageFile.
//
// The expectJSON parameter may be nil, indicating no JSON body.
func CheckMultipart(
	t *testing.T, body io.ReadCloser, h http.Header, v interface{}, expectJSON interface{}, f []api.SendMessageFile,
) {
	_, p, err := mime.ParseMediaType(h.Get("Content-Type"))
	require.NoError(t, err)

	bound, ok := p["boundary"]
	require.True(t, ok, "boundary parameter not set")

	mr := multipart.NewReader(body, bound)

	jsonChecked := false
	// we store the numbers of the missingFiles in a set, so that we know later on,
	// which missingFiles didn't got sent, if any
	missingFiles := make(map[int]struct{}, len(f))

	for i := range f {
		missingFiles[i] = struct{}{}
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		name := part.FormName()

		switch {
		case name == "payload_json":
			checkJSON(t, part, v, expectJSON)

			jsonChecked = true
		case strings.HasPrefix(name, "file"):
			no, err := strconv.Atoi(strings.TrimLeft(name, "file"))
			require.NoErrorf(t, err, `unexpected part with name "%s"`, name)

			if !assert.Lessf(t, no, len(f), "reading file %d, but expected only %d missingFiles", no, len(f)) {
				break
			}

			assert.Equal(t, f[no].Name, part.FileName(), "unequal file names")

			err = equalReader(f[no].Reader, part)
			assert.NoErrorf(t, err, "file %d is not equal to received file", no)

			delete(missingFiles, no)
		default:
			assert.Fail(t, fmt.Sprintf(`unexpected part with name "%s"`, name))
		}
	}

	require.NoError(t, body.Close())

	if !jsonChecked && expectJSON != nil {
		assert.Fail(t, "no json_payload was received, although it was expected")
	}

	if len(missingFiles) > 0 {
		s := joinIntSet(missingFiles, ", ")

		assert.Fail(t, fmt.Sprintf("the files %s did not get sent", s))
	}
}

// CheckQuery checks if the passed query contains the values found in except.
func CheckQuery(t *testing.T, query url.Values, expect url.Values) {
	for name, vals := range query {
		if len(vals) == 0 {
			continue
		}

		expectVal, ok := expect[name]
		if !assert.True(t, ok, "unexpected query field: '"+name+"' with value '"+vals[0]+"'") {
			continue
		}

		assert.Equal(t, expectVal, vals, "query fields for '"+name+"' don't match")

		delete(expect, name)
	}

	for name := range expect {
		assert.Fail(t, "missing query field: '"+name+"'")
	}
}

// checkJSON compares the JSON data from the passed io.Reader against the
// passed expected value.
// v will be used to decode into and should not contain any data.
func checkJSON(t *testing.T, r io.Reader, v interface{}, expect interface{}) {
	err := json.NewDecoder(r).Decode(v)
	require.NoError(t, err)

	val := reflect.ValueOf(expect)
	replaceNullables(val)

	assert.Equal(t, expect, v)
}

// replacesNullables replaces the values of all nullable types with nil, if
// they have assumed their JSON value.
func replaceNullables(val reflect.Value) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)

			t := field.Type()

			elem := field
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}

			isNil := field.Kind() == reflect.Ptr && field.IsNil()

			const initField = "Init"

			// this is a workaround to compensate for json.Unmarshal not calling Unmarshaler
			// functions on JSON null
			switch {
			case t.AssignableTo(nullableBool):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilBool)
				}
			case t.AssignableTo(nullableUint):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilUint)
				}
			case t.AssignableTo(nullableInt):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilInt)
				}
			case t.AssignableTo(nullableString):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilString)
				}
			case t.AssignableTo(nullableColor):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilColor)
				}
			case !isNil && elem.Kind() == reflect.Struct:
				replaceNullables(field)
			}
		}
	}
}

// equalReader checks if the values of the two readers are the same.
func equalReader(r1, r2 io.Reader) error {
	const size = 16

	b1 := make([]byte, size)
	b2 := make([]byte, size)

	for i := 1; ; i++ {
		_, err1 := r1.Read(b1)
		_, err2 := r2.Read(b2)

		if !bytes.Equal(b1, b2) {
			return fmt.Errorf("%d. chunk is not equal:\n%v\nvs.\n%v", i, b1, b2)
		}

		switch {
		case err1 == io.EOF && err2 == io.EOF:
			return nil
		case err1 == io.EOF:
			_, err2 = r2.Read(b2)
			if err2 == io.EOF {
				return nil
			}

			return errors.New("reader 1's stream ended unexpectedly")
		case err2 == io.EOF:
			_, err1 = r1.Read(b1)
			if err1 == io.EOF {
				return nil
			}

			return errors.New("reader 2's stream ended unexpectedly")
		case err1 != nil:
			return err1
		case err2 != nil:
			return err2
		}
	}
}

// strings.Join, but for sets of int.
func joinIntSet(set map[int]struct{}, delim string) string {
	var s string

	first := true

	for no := range set {
		if !first {
			s += delim
		}

		s += strconv.Itoa(no)

		first = false
	}

	return s
}
