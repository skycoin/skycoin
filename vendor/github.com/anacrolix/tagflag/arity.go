package tagflag

import (
	"fmt"
	"reflect"
)

const infArity = 1000

type arity struct {
	min, max int
}

func fieldArity(v reflect.Value, sf reflect.StructField) (arity arity) {
	arity.min = 1
	arity.max = 1
	if v.Kind() == reflect.Slice {
		arity.max = infArity
	}
	if sf.Tag.Get("arity") != "" {
		switch sf.Tag.Get("arity") {
		case "?":
			arity.min = 0
		case "*":
			arity.min = 0
			arity.max = infArity
		case "+":
			arity.max = infArity
		default:
			panic(fmt.Sprintf("unhandled arity tag: %q", sf.Tag.Get("arity")))
		}
	}
	return
}
