package encoder

import (
	"reflect"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
)

type testDynamicHref struct {
	Schema cipher.SHA256
	ObjKey cipher.SHA256
}

type testContainer struct {
	Some string
	Dh   testDynamicHref
	Any  string
}

type testSchema struct {
	Name   string
	Fields []StructField
}

func getSchema(i interface{}) (s testSchema) {
	var (
		typ reflect.Type
		nf  int
	)
	typ = reflect.Indirect(reflect.ValueOf(i)).Type()
	nf = typ.NumField()
	s.Name = typ.Name()
	if nf == 0 {
		return
	}
	s.Fields = make([]StructField, 0, nf)
	for i := 0; i < nf; i++ {
		ft := typ.Field(i)
		if ft.Tag.Get("enc") == "-" || ft.Name == "_" || ft.PkgPath != "" {
			continue
		}
		s.Fields = append(s.Fields, getField(ft))
	}
	return
}

func getField(ft reflect.StructField) (sf StructField) {
	sf.Name = ft.Name
	sf.Type = ft.Type.Name()
	sf.Tag = string(ft.Tag)
	sf.Kind = uint32(ft.Type.Kind())
	return
}

// Actually it's hard or impossible with current implementation
// the encoder treat
//
//
func Test_fields(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Error("unexpected panicing")
		}
	}()
	c := testContainer{
		Some: "some",
		Dh: testDynamicHref{
			cipher.SHA256{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			cipher.SHA256{0, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		Any: "any",
	}
	s := getSchema(testContainer{})
	t.Logf("Schema %v", s)
	in := Serialize(c)
	// type name
	var dhTypeName = reflect.TypeOf(testDynamicHref{}).Name()
	t.Log("dhTypeName: ", dhTypeName)
	for _, sf := range s.Fields {
		if sf.Type == dhTypeName {
			t.Log("Field ", sf)
			var dh testDynamicHref
			if err := DeserializeField(in, s.Fields, sf.Name, &dh); err != nil {
				t.Error("unexpected error:", err)
			}
			if dh.Schema != (cipher.SHA256{1, 2, 3}) {
				t.Error("wrong schema of DH")
			}
			if dh.ObjKey != (cipher.SHA256{3, 4, 5}) {
				t.Error("wrong object key of DH")
			}
		} else {
			t.Log("Field ", sf)
			var st string
			if err := DeserializeField(in, s.Fields, sf.Name, &st); err != nil {
				t.Error("unexpected error:", err)
			}
			switch st {
			case "some", "any":
			default:
				t.Errorf("unexpected %s value %q", sf.Name, st)
			}
		}
	}
}
