package encoder

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

type ReflectStructure struct {
	Name   string
	Fields []ReflectionField
}
type ReflectionField struct {
	Name string
	Type string
}

func FieldData(data interface{}) (ReflectStructure, error) {
	var err error
	value := reflect.Indirect(reflect.ValueOf(data))
	ref, err := getFieldType(value)

	//need to make json of this

	// line := fmt.Sprintf("%v", fieldType)
	// lines := strings.Split(line, " ")
	// lines = append(lines, line)
	// if err := writeLines(lines, "file.txt"); err != nil {
	// }

	return ref, err
}

func getFieldType(v reflect.Value) (ReflectStructure, error) {
	v = reflect.Indirect(v)
	var err error
	var result ReflectStructure
	var fields []ReflectionField
	typeOfT := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		name := typeOfT.Field(i).Name
		fieldType := fmt.Sprint("", f.Kind())

		values := ReflectionField{
			Name: name,
			Type: fieldType,
		}
		fields = append(fields, values)
	}
	result.Name = "ClonedStruct"

	result.Fields = fields
	return result, err
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
