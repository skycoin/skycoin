package encoder

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
)

// StructField field struct
type StructField struct {
	Name string `json:"name"`
	Kind uint32 `json:"kind"`
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

func (s *StructField) String() string {
	return fmt.Sprintln(s.Name, s.Type, s.Tag)
}

//TODO: replace fieldType on reflect.Kind
func getFieldSize(in []byte, d *decoder, fieldType reflect.Kind, s int) (int, error) {
	switch fieldType {
	case reflect.Slice, reflect.String:
		length := int(leUint32(d.buf[s : s+4]))
		s += 4 + length
	case reflect.Struct, reflect.Array:
		s += 32
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		s++
	case reflect.Int16, reflect.Uint16:
		s += 2
	case reflect.Int32, reflect.Uint32:
		s += 4
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		s += 8
	default:
		fmt.Println(fieldType)
		log.Panicf("Decode error: kind %s not handled", fieldType)
	}
	return s, nil
}

//TODO: replace fieldType on reflect.Kind
func getFieldValue(in []byte, d *decoder, fieldType reflect.Kind, s int) string {
	fd := &decoder{buf: make([]byte, len(in)-s)}
	copy(fd.buf, d.buf[s:])
	switch fieldType {
	case reflect.Slice, reflect.String:
		length := int(leUint32(fd.buf[0:4]))
		return string(fd.buf[4 : 4+length])
	case reflect.Struct, reflect.Array:
		s := cipher.SHA256{}
		s.Set(fd.buf[0:32])
		return s.Hex()
	case reflect.Bool:
		return strconv.FormatBool(fd.bool())
	case reflect.Int8:
		return strconv.Itoa(int(fd.int8()))
	case reflect.Int16:
		return strconv.Itoa(int(fd.int16()))
	case reflect.Int32:
		return strconv.Itoa(int(fd.int32()))
	case reflect.Int64:
		return strconv.Itoa(int(fd.int64()))
	case reflect.Uint8:
		return strconv.Itoa(int(fd.uint8()))
	case reflect.Uint16:
		return strconv.Itoa(int(fd.uint16()))
	case reflect.Uint32:
		return strconv.Itoa(int(fd.uint32()))
	case reflect.Uint64:
		return strconv.Itoa(int(fd.uint64()))
	default:
		log.Panicf("Decode error: kind %s not handled", fieldType)
	}
	return ""
}

// DeserializeField deserialize field
func DeserializeField(in []byte, fields []StructField, fieldName string, field interface{}) error {

	d := &decoder{buf: make([]byte, len(in))}
	copy(d.buf, in)
	fv := reflect.ValueOf(field).Elem()
	s := 0
	for _, f := range fields {
		if f.Name == fieldName {
			fd := &decoder{buf: make([]byte, len(in)-s)}
			copy(fd.buf, d.buf[s:])
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

// ParseFields parse fields
func ParseFields(in []byte, fields []StructField) map[string]string {
	result := map[string]string{}
	d := &decoder{buf: make([]byte, len(in))}
	copy(d.buf, in)
	s := 0
	for _, f := range fields {
		resShift, _ := getFieldSize(in, d, reflect.Kind(f.Kind), s)
		result[f.Name] = getFieldValue(in, d, reflect.Kind(f.Kind), s)
		s = resShift
	}
	return result
}
