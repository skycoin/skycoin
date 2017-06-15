// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Marshalling and unmarshalling of
// bit torrent bencode data into Go structs using reflection.
//
// Based upon the standard Go language JSON package.

package bencode

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

type structBuilder struct {
	val reflect.Value

	// if map_ != nil, write val to map_[key] on each change
	map_ reflect.Value
	key  reflect.Value
}

var nobuilder *structBuilder

func isfloat(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func setfloat(v reflect.Value, f float64) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		v.SetFloat(f)
	}
}

func setint(val reflect.Value, i int64) {
	switch v := val; v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v.SetUint(uint64(i))
	case reflect.Interface:
		v.Set(reflect.ValueOf(i))
	default:
		panic("setint called for bogus type: " + val.Kind().String())
	}
}

// If updating b.val is not enough to update the original,
// copy a changed b.val out to the original.
func (b *structBuilder) Flush() {
	if b == nil {
		return
	}
	if b.map_.IsValid() {
		b.map_.SetMapIndex(b.key, b.val)
	}
}

func (b *structBuilder) Int64(i int64) {
	if b == nil {
		return
	}
	if !b.val.CanSet() {
		x := 0
		b.val = reflect.ValueOf(&x).Elem()
	}
	v := b.val
	if isfloat(v) {
		setfloat(v, float64(i))
	} else {
		setint(v, i)
	}
}

func (b *structBuilder) Uint64(i uint64) {
	if b == nil {
		return
	}
	if !b.val.CanSet() {
		x := uint64(0)
		b.val = reflect.ValueOf(&x).Elem()
	}
	v := b.val
	if isfloat(v) {
		setfloat(v, float64(i))
	} else {
		setint(v, int64(i))
	}
}

func (b *structBuilder) Float64(f float64) {
	if b == nil {
		return
	}
	if !b.val.CanSet() {
		x := float64(0)
		b.val = reflect.ValueOf(&x).Elem()
	}
	v := b.val
	if isfloat(v) {
		setfloat(v, f)
	} else {
		setint(v, int64(f))
	}
}

func (b *structBuilder) String(s string) {
	if b == nil {
		return
	}

	switch b.val.Kind() {
	case reflect.String:
		if !b.val.CanSet() {
			x := ""
			b.val = reflect.ValueOf(&x).Elem()

		}
		b.val.SetString(s)
	case reflect.Interface:
		b.val.Set(reflect.ValueOf(s))
	}
}

func (b *structBuilder) Array() {
	if b == nil {
		return
	}
	if v := b.val; v.Kind() == reflect.Slice {
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 8))
		}
	}
}

func (b *structBuilder) Elem(i int) builder {
	if b == nil || i < 0 {
		return nobuilder
	}
	switch v := b.val; v.Kind() {
	case reflect.Array:
		if i < v.Len() {
			return &structBuilder{val: v.Index(i)}
		}
	case reflect.Slice:
		if i >= v.Cap() {
			n := v.Cap()
			if n < 8 {
				n = 8
			}
			for n <= i {
				n *= 2
			}
			nv := reflect.MakeSlice(v.Type(), v.Len(), n)
			reflect.Copy(nv, v)
			v.Set(nv)
		}
		if v.Len() <= i && i < v.Cap() {
			v.SetLen(i + 1)
		}
		if i < v.Len() {
			return &structBuilder{val: v.Index(i)}
		}
	}
	return nobuilder
}

func (b *structBuilder) Map() {
	if b == nil {
		return
	}
	if v := b.val; v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.Zero(v.Type().Elem()).Addr())
			b.Flush()
		}
		b.map_ = reflect.Value{}
		b.val = v.Elem()
	}
	if v := b.val; v.Kind() == reflect.Map && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
}

func (b *structBuilder) Key(k string) builder {
	if b == nil {
		return nobuilder
	}
	switch v := reflect.Indirect(b.val); v.Kind() {
	case reflect.Struct:
		t := v.Type()
		// Case-insensitive field lookup.
		k = strings.ToLower(k)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			key := bencodeKey(field, nil)
			if strings.ToLower(key) == k ||
				strings.ToLower(field.Name) == k {
				return &structBuilder{val: v.Field(i)}
			}
		}
	case reflect.Map:
		t := v.Type()
		if t.Key() != reflect.TypeOf(k) {
			break
		}
		key := reflect.ValueOf(k)
		elem := v.MapIndex(key)
		if !elem.IsValid() {
			v.SetMapIndex(key, reflect.Zero(t.Elem()))
			elem = v.MapIndex(key)
		}
		return &structBuilder{val: elem, map_: v, key: key}
	}
	return nobuilder
}

// Unmarshal reads and parses the bencode syntax data from r and fills in
// an arbitrary struct or slice pointed at by val.
// It uses the reflect package to assign to fields
// and arrays embedded in val.  Well-formed data that does not fit
// into the struct is discarded.
//
// For example, given these definitions:
//
//	type Email struct {
//		Where string;
//		Addr string;
//	}
//
//	type Result struct {
//		Name string;
//		Phone string;
//		Email []Email
//	}
//
//	var r = Result{ "name", "phone", nil }
//
// unmarshalling the bencode syntax string
//
//	"d5:emailld5:where4:home4:addr15:gre@example.come\
//  d5:where4:work4:addr12:gre@work.comee4:name14:Gr\
//  ace R. Emlin7:address15:123 Main Streete"
//
// via Unmarshal(s, &r) is equivalent to assigning
//
//	r = Result{
//		"Grace R. Emlin",	// name
//		"phone",		// no phone given
//		[]Email{
//			Email{ "home", "gre@example.com" },
//			Email{ "work", "gre@work.com" }
//		}
//	}
//
// Note that the field r.Phone has not been modified and
// that the bencode field "address" was discarded.
//
// Because Unmarshal uses the reflect package, it can only
// assign to upper case fields.  Unmarshal uses a case-insensitive
// comparison to match bencode field names to struct field names.
//
// If you provide a tag string for a struct member, the tag string
// will be used as the bencode dictionary key for that member.
// Bencode undestands both the original single-string and updated
// list-of-key-value-pairs tag string syntax. The list-of-key-value
// pairs syntax is assumed, with a fallback to the original single-string
// syntax. The key for bencode values is bencode.
//
// To unmarshal a top-level bencode array, pass in a pointer to an empty
// slice of the correct type.
//
func Unmarshal(r io.Reader, val interface{}) (err error) {
	// If e represents a value, the answer won't get back to the
	// caller.  Make sure it's a pointer.
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		err = errors.New("Attempt to unmarshal into a non-pointer")
		return
	}
	err = unmarshalValue(r, reflect.Indirect(reflect.ValueOf(val)))
	return
}

func unmarshalValue(r io.Reader, v reflect.Value) (err error) {
	var b *structBuilder

	// XXX: Decide if the extra codnitions are needed. Affect map?
	if ptr := v; ptr.Kind() == reflect.Ptr {
		if slice := ptr.Elem(); slice.Kind() == reflect.Slice || slice.Kind() == reflect.Int || slice.Kind() == reflect.String {
			b = &structBuilder{val: slice}
		}
	}

	if b == nil {
		b = &structBuilder{val: v}
	}

	err = parse(r, b)
	return
}

type MarshalError struct {
	T reflect.Type
}

func (e *MarshalError) Error() string {
	return "bencode cannot encode value of type " + e.T.String()
}

func writeArrayOrSlice(w io.Writer, val reflect.Value) (err error) {
	_, err = fmt.Fprint(w, "l")
	if err != nil {
		return
	}
	for i := 0; i < val.Len(); i++ {
		if err := writeValue(w, val.Index(i)); err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "e")
	if err != nil {
		return
	}
	return nil
}

type stringValue struct {
	key       string
	value     reflect.Value
	omitEmpty bool
}

type stringValueArray []stringValue

// Satisfy sort.Interface

func (a stringValueArray) Len() int { return len(a) }

func (a stringValueArray) Less(i, j int) bool { return a[i].key < a[j].key }

func (a stringValueArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func writeSVList(w io.Writer, svList stringValueArray) (err error) {
	sort.Sort(svList)

	for _, sv := range svList {
		if sv.isValueNil() {
			continue // Skip null values
		}
		s := sv.key
		_, err = fmt.Fprintf(w, "%d:%s", len(s), s)
		if err != nil {
			return
		}

		if err = writeValue(w, sv.value); err != nil {
			return
		}
	}
	return
}

func writeMap(w io.Writer, val reflect.Value) (err error) {
	key := val.Type().Key()
	if key.Kind() != reflect.String {
		return &MarshalError{val.Type()}
	}
	_, err = fmt.Fprint(w, "d")
	if err != nil {
		return
	}

	keys := val.MapKeys()

	// Sort keys

	svList := make(stringValueArray, len(keys))
	for i, key := range keys {
		svList[i].key = key.String()
		svList[i].value = val.MapIndex(key)
	}

	err = writeSVList(w, svList)
	if err != nil {
		return
	}

	_, err = fmt.Fprint(w, "e")
	if err != nil {
		return
	}
	return
}

func bencodeKey(field reflect.StructField, sv *stringValue) (key string) {
	key = field.Name
	tag := field.Tag
	if len(tag) > 0 {
		// Backwards compatability
		// If there's a bencode key/value entry in the tag, use it.
		var tagOpt tagOptions
		key, tagOpt = parseTag(tag.Get("bencode"))
		if len(key) == 0 {
			key = tag.Get("bencode")
			if len(key) == 0 && !strings.Contains(string(tag), ":") {
				// If there is no ":" in the tag, assume it is an old-style tag.
				key = string(tag)
			} else {
				key = field.Name
			}
		}
		if sv != nil && tagOpt.Contains("omitempty") {
			sv.omitEmpty = true
		}
	}
	if sv != nil {
		sv.key = key
	}
	return
}

// tagOptions is the string following a comma in a struct field's "bencode"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's bencode tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

func writeStruct(w io.Writer, val reflect.Value) (err error) {
	_, err = fmt.Fprint(w, "d")
	if err != nil {
		return
	}

	typ := val.Type()

	numFields := val.NumField()
	svList := make(stringValueArray, numFields)

	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		bencodeKey(field, &svList[i])
		// The tag `bencode:"-"` should mean that this field must be ignored
		// See https://golang.org/pkg/encoding/json/#Marshal or https://golang.org/pkg/encoding/xml/#Marshal
		// We set a zero value so that it is ignored by the writeSVList() function
		if svList[i].key == "-" {
			svList[i].value = reflect.Value{}
		} else {
			svList[i].value = val.Field(i)
		}
	}

	err = writeSVList(w, svList)
	if err != nil {
		return
	}

	_, err = fmt.Fprint(w, "e")
	if err != nil {
		return
	}
	return
}

func writeValue(w io.Writer, val reflect.Value) (err error) {
	if !val.IsValid() {
		err = errors.New("Can't write null value")
		return
	}

	switch v := val; v.Kind() {
	case reflect.String:
		s := v.String()
		_, err = fmt.Fprintf(w, "%d:%s", len(s), s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err = fmt.Fprintf(w, "i%de", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		_, err = fmt.Fprintf(w, "i%de", v.Uint())
	case reflect.Array:
		err = writeArrayOrSlice(w, v)
	case reflect.Slice:
		switch val.Type().String() {
		case "[]uint8":
			// special case as byte-string
			s := string(v.Bytes())
			_, err = fmt.Fprintf(w, "%d:%s", len(s), s)
		default:
			err = writeArrayOrSlice(w, v)
		}
	case reflect.Map:
		err = writeMap(w, v)
	case reflect.Struct:
		err = writeStruct(w, v)
	case reflect.Interface:
		err = writeValue(w, v.Elem())
	default:
		err = &MarshalError{val.Type()}
	}
	return
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (sv stringValue) isValueNil() bool {
	if !sv.value.IsValid() || (sv.omitEmpty && isEmptyValue(sv.value)) {
		return true
	}
	switch v := sv.value; v.Kind() {
	case reflect.Interface:
		return !v.Elem().IsValid()
	}
	return false
}

// Marshal writes the bencode encoding of val to w.
//
// Marshal traverses the value v recursively.
//
// Marshal uses the following type-dependent encodings:
//
// Floating point, integer, and Number values encode as bencode numbers.
//
// String values encode as bencode strings.
//
// Array and slice values encode as bencode arrays.
//
// Struct values encode as bencode maps. Each exported struct field
// becomes a member of the object.
// The object's default key string is the struct field name
// but can be specified in the struct field's tag value. The text of
// the struct field's tag value is the key name. Examples:
//
//   // Field appears in bencode as key "Field".
//   Field int
//
//   // Field appears in bencode as key "myName".
//   Field int "myName"
//
// Anonymous struct fields are ignored.
//
// Map values encode as bencode objects.
// The map's key type must be string; the object keys are used directly
// as map keys.
//
// Boolean, Pointer, Interface, Channel, complex, and function values cannot
// be encoded in bencode.
// Attempting to encode such a value causes Marshal to return
// a MarshalError.
//
// Bencode cannot represent cyclic data structures and Marshal does not
// handle them.  Passing cyclic structures to Marshal will result in
// an infinite recursion.
//
func Marshal(w io.Writer, val interface{}) error {
	return writeValue(w, reflect.ValueOf(val))
}
