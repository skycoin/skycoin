package encoder

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
)

type ReflectionField struct {
	Name string
	Type string
}

type Sha254Elements struct {
	sha []string
}

func FieldData(data interface{}) (Sha254Elements, error) {
	var err error
	value := reflect.Indirect(reflect.ValueOf(data))
	ref, err := getFieldType(value)
	var result Sha254Elements

	for _, rf := range ref {
		bv := []byte(rf.Name + rf.Type)
		hasher := sha256.New()
		hasher.Write(bv)
		shaElement := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
		result.sha = append(result.sha, shaElement)
	}

	return result, err
}

func getFieldType(v reflect.Value) ([]ReflectionField, error) {
	v = reflect.Indirect(v)
	var err error
	var fields []ReflectionField
	typeOfT := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		if f.Kind() == reflect.Struct {
			fields, err = getFieldType(f)
			if err != nil {
			}
		}

		name := typeOfT.Field(i).Name
		fieldType := fmt.Sprint("", f.Kind())

		values := ReflectionField{
			Name: name,
			Type: fieldType,
		}
		fields = append(fields, values)
	}

	return fields, err
}

func writeLines(lines []string, path string) error {
	// line := fmt.Sprintf("%v", fieldType)
	// lines := strings.Split(line, " ")
	// lines = append(lines, line)
	// if err := writeLines(lines, "filename.txt"); err != nil {
	// }

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
