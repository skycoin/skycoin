package tagflag

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/bradfitz/iter"
)

const flagPrefix = "-"

// Walks the fields of the given struct, calling the function with the value
// and StructField for each field. Returning true from the function will halt
// traversal.
func foreachStructField(_struct reflect.Value, f func(fv reflect.Value, sf reflect.StructField) (stop bool)) {
	t := _struct.Type()
	for i := range iter.N(t.NumField()) {
		sf := t.Field(i)
		fv := _struct.Field(i)
		if f(fv, sf) {
			break
		}
	}
}

func canMarshal(f reflect.Value) bool {
	return valueMarshaler(f) != nil
}

// Returns a marshaler for the given value, or nil if there isn't one.
func valueMarshaler(v reflect.Value) marshaler {
	if v.CanAddr() {
		if am, ok := v.Addr().Interface().(Marshaler); ok {
			return dynamicMarshaler{
				marshal:               func(_ reflect.Value, s string) error { return am.Marshal(s) },
				explicitValueRequired: am.RequiresExplicitValue(),
			}
		}
	}
	if bm, ok := builtinMarshalers[v.Type()]; ok {
		return bm
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Struct:
		return nil
	case reflect.Bool:
		return dynamicMarshaler{
			marshal: func(v reflect.Value, s string) error {
				if s == "" {
					v.SetBool(true)
					return nil
				}
				b, err := strconv.ParseBool(s)
				v.SetBool(b)
				return err
			},
			explicitValueRequired: false,
		}
	}
	return defaultMarshaler{}
}

// Turn a struct field name into a flag name. In particular this lower cases
// leading acronyms, and the first capital letter.
func fieldFlagName(fieldName string) (ret string) {
	// defer func() { log.Println(fieldName, ret) }()
	// TCP
	if ss := regexp.MustCompile("^[[:upper:]]{2,}$").FindStringSubmatch(fieldName); ss != nil {
		return strings.ToLower(ss[0])
	}
	// TCPAddr
	if ss := regexp.MustCompile("^([[:upper:]]+)([[:upper:]][^[:upper:]].*?)$").FindStringSubmatch(fieldName); ss != nil {
		return strings.ToLower(ss[1]) + ss[2]
	}
	// Addr
	if ss := regexp.MustCompile("^([[:upper:]])(.*)$").FindStringSubmatch(fieldName); ss != nil {
		return strings.ToLower(ss[1]) + ss[2]
	}
	panic(fieldName)
}
