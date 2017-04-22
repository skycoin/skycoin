package encoder

import (
	"fmt"
	"reflect"

	"github.com/skycoin/skycoin/src/cipher"
)

type StructField struct {
	Name string `json:"name"`
	Kind uint32 `json:"kind"`
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

func (s *StructField) String() string {
	return fmt.Sprintln(s.Name, s.Type, s.Tag)
}

func getFieldSize(in []byte, d *decoder, fieldType reflect.Kind,
	s int) (int, error) {

	switch fieldType {
	case reflect.Slice, reflect.String:
		length := int(le_Uint32(d.buf[s : s+4]))
		s += 4 + length
	case reflect.Struct, reflect.Array:
		s += 32
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		s += 1
	case reflect.Int16, reflect.Uint16:
		s += 2
	case reflect.Int32, reflect.Uint32:
		s += 4
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		s += 8
	default:
		return 0, fmt.Errorf("Decode error: kind %s not handled", fieldType)
	}
	return s, nil
}

func getFieldValue(in []byte, d *decoder, fieldType reflect.Kind,
	s int) (v interface{}, err error) {

	fd := &decoder{buf: make([]byte, len(in)-s)}
	copy(fd.buf, d.buf[s:])
	switch fieldType {
	case reflect.Slice, reflect.String:
		length := int(le_Uint32(fd.buf[0:4]))
		v = string(fd.buf[4 : 4+length])
	case reflect.Struct, reflect.Array:
		s := cipher.SHA256{}
		s.Set(fd.buf[0:32])
		v = s.Hex()
	case reflect.Bool:
		v = fd.bool()
	case reflect.Int8:
		v = fd.int8()
	case reflect.Int16:
		v = fd.int16()
	case reflect.Int32:
		v = fd.int32()
	case reflect.Int64:
		v = fd.int64()
	case reflect.Uint8:
		v = fd.uint8()
	case reflect.Uint16:
		v = fd.uint16()
	case reflect.Uint32:
		v = fd.uint32()
	case reflect.Uint64:
		v = fd.uint64()
	default:
		err = fmt.Errorf("Decode error: kind %s not handled", fieldType)
	}
	return
}

func DeserializeField(in []byte, fields []StructField, fieldName string,
	field interface{}) error {

	d := &decoder{buf: in}
	fv := reflect.ValueOf(field).Elem()
	s := 0
	for _, f := range fields {
		if f.Name == fieldName {
			fd := &decoder{buf: d.buf[s:]}
			fd.value(fv)
			return nil
		}
		res, err := getFieldSize(in, d, reflect.Kind(f.Kind), s)
		if err != nil {
			return err
		}
		s = res
	}
	return nil
}

func ParseFields(in []byte,
	fields []StructField) (msi map[string]interface{}, err error) {

	var (
		d     *decoder = &decoder{buf: in}
		s     int      = 0
		shift int
	)

	msi = map[string]interface{}{}

	for _, f := range fields {
		shift, err = getFieldSize(in, d, reflect.Kind(f.Kind), s)
		if err != nil {
			return
		}
		msi[f.Name], err = getFieldValue(in, d, reflect.Kind(f.Kind), s)
		if err != nil {
			return
		}
		s = shift
	}
	return
}
